package mrtypes

import (
	"context"
	"time"
)

type ModemRebooter interface {
	Reboot(ctx context.Context, conf Config) error
}

type Config struct {
	Type          string `json:"type"`
	Address       string `json:"address"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
}

type RebootConfig struct {
	RebootAfterDownFor time.Duration
	ModemRecoversIn    time.Duration
}
