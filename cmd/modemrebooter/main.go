package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/function61/gokit/logger"
	"github.com/function61/gokit/ossignal"
	"github.com/function61/gokit/systemdinstaller"
	"github.com/joonas-fi/modemrebooter/pkg/internetupdetector"
	"github.com/joonas-fi/modemrebooter/pkg/mrtypes"
	"github.com/joonas-fi/modemrebooter/pkg/tplinktlmr6400"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var version = "dev" // replaced dynamically at build time

const (
	tagline = "Reboots your modem if internet is down"
)

var defaultRebootConfig = mrtypes.RebootConfig{
	RebootAfterDownFor: 4 * time.Minute,
	ModemRecoversIn:    4 * time.Minute,
}

func mainLoop(ctx context.Context, conf mrtypes.Config) error {
	log := logger.New("internet")

	rebooter, err := initRebooter(conf)
	if err != nil {
		return err
	}

	state := State{}

	for {
		up := internetupdetector.IsInternetUp(ctx)

		previousState := state

		if up {
			state = state.Up()
		} else {
			state = state.Down(time.Now())
		}

		if state.IsUpDifferentTo(previousState) {
			if up {
				log.Info("came back UP")
			} else {
				log.Error("went DOWN")
			}
		}

		if up {
			log.Debug("up")
		} else {
			log.Info(fmt.Sprintf("down for %s", time.Since(state.wentDownAt)))
		}

		if state.ShouldReboot(defaultRebootConfig, time.Now()) {
			log.Info("rebooting modem")

			if err := rebooter.Reboot(conf); err != nil {
				log.Error(fmt.Sprintf("reboot failed: %s", err.Error()))
			} else {
				log.Info("reboot succeeded")

				state = state.SuccesfullReboot(time.Now())
			}
		}

		select {
		case <-ctx.Done():
			return nil // graceful stop
		case <-time.After(1 * time.Minute):
		}
	}
}

func main() {
	app := &cobra.Command{
		Use:     os.Args[0],
		Short:   tagline,
		Version: version,
	}

	app.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Runs the program",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New("main")
			log.Info(fmt.Sprintf("starting %s", version))
			defer log.Info("stopped")

			conf, err := readConfig()
			if err != nil {
				panic(err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				log.Info(fmt.Sprintf("got %s; stopping", ossignal.WaitForInterruptOrTerminate()))
				cancel()
			}()

			if err := mainLoop(ctx, *conf); err != nil {
				panic(err)
			}
		},
	})
	app.AddCommand(writeSystemdFileEntry())

	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeSystemdFileEntry() *cobra.Command {
	return &cobra.Command{
		Use:   "write-systemd-unit-file",
		Short: "Install unit file to start this application on startup",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			systemdHints, err := systemdinstaller.InstallSystemdServiceFile("modemrebooter", []string{"run"}, tagline)
			if err != nil {
				panic(err)
			}

			fmt.Println(systemdHints)
		},
	}
}

func readConfig() (*mrtypes.Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	jsonDecoder.DisallowUnknownFields()

	conf := &mrtypes.Config{}
	if err := jsonDecoder.Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func initRebooter(conf mrtypes.Config) (mrtypes.ModemRebooter, error) {
	switch conf.Type {
	case "tplinktlmr6400":
		return tplinktlmr6400.New(), nil
	default:
		return nil, fmt.Errorf("unknown modem type: %s", conf.Type)
	}
}
