package main

import (
	"github.com/function61/gokit/assert"
	"testing"
	"time"
)

func TestState(t *testing.T) {
	state := State{}

	midnight := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tplus := func(minutes time.Duration) time.Time { return midnight.Add(minutes * time.Minute) }

	// all transition methods create copies of previous state
	state = state.Up()
	state = state.Up()

	// base state is UP
	assert.Assert(t, state.IsUpDifferentTo(state.Up()) == false)                          // up==up?
	assert.Assert(t, state.IsUpDifferentTo(state.Down(midnight)) == true)                 // up==down?
	assert.Assert(t, state.Down(midnight).IsUpDifferentTo(state) == true)                 // down==up?
	assert.Assert(t, state.Down(midnight).IsUpDifferentTo(state.Down(midnight)) == false) // down==down?

	assert.Assert(t, state.wentDownAt.IsZero())
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, midnight))

	state = state.Down(midnight)

	assert.Assert(t, state.wentDownAt.Equal(midnight))

	// reboot should be only possible at 5 minute mark
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(1)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(2)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(3)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(4)))
	assert.Assert(t, state.ShouldReboot(defaultRebootConfig, tplus(5)))

	// now reboot
	state = state.SuccesfullReboot(tplus(5))

	// internet keeps being down, but reboot is not possible until "modemRecoversIn"
	// from last reboot
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(5)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(6)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(7)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(8)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(9)))

	// another reboot after previous reboot
	assert.Assert(t, state.ShouldReboot(defaultRebootConfig, tplus(10)))

	state = state.SuccesfullReboot(tplus(10))

	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(10)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(11)))

	// internet went back UP, woohoo!
	state = state.Up()

	// while we're up, should not reboot
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(11)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(12)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(40)))

	// down again :(
	state = state.Down(tplus(40))

	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(40)))
	assert.Assert(t, !state.ShouldReboot(defaultRebootConfig, tplus(41)))
	assert.Assert(t, state.ShouldReboot(defaultRebootConfig, tplus(45)))
}
