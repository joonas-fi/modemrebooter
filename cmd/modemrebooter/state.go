package main

import (
	"github.com/joonas-fi/modemrebooter/pkg/mrtypes"
	"time"
)

type State struct {
	wentDownAt             time.Time
	lastSuccesfullRebootAt time.Time
}

func (s State) IsUpDifferentTo(other State) bool {
	return !s.wentDownAt.Equal(other.wentDownAt)
}

func (s State) Up() State {
	// reset possible wentDownAt to zero
	return State{
		lastSuccesfullRebootAt: s.lastSuccesfullRebootAt,
	}
}

func (s State) Down(now time.Time) State {
	if s.wentDownAt.IsZero() { // wasn't down before => set wentDownAt to now
		return State{
			lastSuccesfullRebootAt: s.lastSuccesfullRebootAt,
			wentDownAt:             now,
		}
	}

	return s
}

func (s State) SuccesfullReboot(now time.Time) State {
	// not resetting wentDownAt, because succesfull reboot does not
	// guarantee internet will be back up
	return State{
		lastSuccesfullRebootAt: now,
		wentDownAt:             s.wentDownAt,
	}
}

func (s State) ShouldReboot(rc mrtypes.RebootConfig, now time.Time) bool {
	return !s.wentDownAt.IsZero() &&
		now.Sub(s.wentDownAt) > rc.RebootAfterDownFor &&
		now.Sub(s.lastSuccesfullRebootAt) > rc.ModemRecoversIn
}
