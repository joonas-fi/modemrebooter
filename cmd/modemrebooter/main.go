package main

import (
	"context"
	"fmt"
	"github.com/function61/gokit/dynversion"
	"github.com/function61/gokit/jsonfile"
	"github.com/function61/gokit/logex"
	"github.com/function61/gokit/ossignal"
	"github.com/function61/gokit/systemdinstaller"
	"github.com/joonas-fi/modemrebooter/pkg/internetupdetector"
	"github.com/joonas-fi/modemrebooter/pkg/mrtypes"
	"github.com/joonas-fi/modemrebooter/pkg/tplinktlmr6400"
	"github.com/joonas-fi/modemrebooter/pkg/zyxelvmg1312b10d"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	tagline = "Reboots your modem if internet is down"
)

var defaultRebootConfig = mrtypes.RebootConfig{
	RebootAfterDownFor: 4 * time.Minute,
	ModemRecoversIn:    4 * time.Minute,
}

func mainLoop(ctx context.Context, conf mrtypes.Config, logl *logex.Leveled) error {
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
				logl.Info.Println("came back UP")
			} else {
				logl.Error.Println("went DOWN")
			}
		}

		if up {
			logl.Debug.Println("up")
		} else {
			logl.Info.Printf("down for %s", time.Since(state.wentDownAt))
		}

		if state.ShouldReboot(defaultRebootConfig, time.Now()) {
			logl.Info.Println("rebooting modem")

			// modem reboot should succeed within 60 seconds
			rebootCtx, cancel := context.WithTimeout(ctx, 60*time.Second)

			if err := rebooter.Reboot(rebootCtx, conf); err != nil {
				logl.Error.Printf("reboot failed: %s", err.Error())
			} else {
				logl.Info.Println("reboot succeeded")

				state = state.SuccesfullReboot(time.Now())
			}

			cancel()
		}

		select {
		case <-ctx.Done():
			return nil // graceful stop
		case <-time.After(1 * time.Minute):
		}
	}
}

func mainInternal() error {
	rootLogger := logex.StandardLogger()

	mainLogger := logex.Levels(logex.Prefix("main", rootLogger))

	mainLogger.Info.Printf("starting %s", dynversion.Version)
	defer mainLogger.Info.Println("stopped")

	conf, err := readConfig()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		mainLogger.Info.Printf("got %s; stopping", <-ossignal.InterruptOrTerminate())
	}()

	if err := mainLoop(ctx, *conf, logex.Levels(logex.Prefix("internet", rootLogger))); err != nil {
		return err
	}

	return nil
}

func runReboot() error {
	conf, err := readConfig()
	if err != nil {
		return err
	}

	rebooter, err := initRebooter(*conf)
	if err != nil {
		return err
	}

	// modem reboot should succeed within 60 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return rebooter.Reboot(ctx, *conf)
}

func main() {
	app := &cobra.Command{
		Use:     os.Args[0],
		Short:   tagline,
		Version: dynversion.Version,
	}

	app.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Runs the program",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mainInternal(); err != nil {
				panic(err)
			}
		},
	})
	app.AddCommand(&cobra.Command{
		Use:   "reboot",
		Short: "Just reboots the modem now",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runReboot(); err != nil {
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
			service := systemdinstaller.SystemdServiceFile(
				"modemrebooter",
				tagline,
				systemdinstaller.Args("run"),
				systemdinstaller.Docs("https://github.com/joonas-fi/modemrebooter"))

			if err := systemdinstaller.Install(service); err != nil {
				panic(err)
			}

			fmt.Println(systemdinstaller.GetHints(service))
		},
	}
}

func initRebooter(conf mrtypes.Config) (mrtypes.ModemRebooter, error) {
	switch conf.Type {
	case "zyxelvmg1312b10d":
		return zyxelvmg1312b10d.New(), nil
	case "tplinktlmr6400":
		return tplinktlmr6400.New(), nil
	default:
		return nil, fmt.Errorf("unknown modem type: %s", conf.Type)
	}
}

func readConfig() (*mrtypes.Config, error) {
	conf := &mrtypes.Config{}
	return conf, jsonfile.Read("config.json", conf, true)
}
