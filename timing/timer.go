package timing

import (
	"fmt"
	"time"
)

// Timer wraps state of current task or break
type Timer struct {
	IsPaused    bool
	IsActive    bool
	CurrentTodo string
	Elapsed     int
	remaining   int
	Total       int
	IdlingSince time.Time
}

// PauseToggle controls the timer
func (t *Timer) PauseToggle() {
	t.IsPaused = !t.IsPaused
}

// Cancel stops the timer
func (t *Timer) Cancel() {
	t.remaining = 0
}

// Reset all internal state
func (t *Timer) Reset() {
	t.IsPaused = false
	t.IsActive = false
	t.Elapsed = 0
	t.CurrentTodo = ""
}
func (t *Timer) Idle() {
	t.IdlingSince = time.Now()
}

// NewTimer constructor
func NewTimer(remaining int) *Timer {
	return &Timer{remaining: remaining, Total: remaining}
}

// Progress proceeds the timer 1 second
func (t *Timer) Progress() bool {
	if t.remaining == 0 {
		t.IsActive = false
		return false
	}
	for t.IsPaused {
		time.Sleep(time.Millisecond * 100)
	}
	t.remaining--
	t.Elapsed++
	t.IsActive = true

	return true
}

func (t *Timer) IdlingTime() time.Duration {
	return time.Since(t.IdlingSince)
}

// RemainingToString returns remaining time in mm:ss
func (t *Timer) RemainingToString() string {
	return fmt.Sprintf("%02d:%02d", t.remaining/60, t.remaining%60)
}

// ElapsedToString returns formatted elapsed time
func (t *Timer) ElapsedToString() string {
	if t.Elapsed/3600 > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", t.Elapsed/3600, t.Elapsed/60, t.Elapsed%60)
	}

	return fmt.Sprintf("%02d:%02d", t.Elapsed/60, t.Elapsed%60)
}
