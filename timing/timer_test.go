package timing

import (
	"testing"

	"github.com/stvp/assert"
)

func TestTimer_PauseToggle(t *testing.T) {
	timer := &Timer{IsPaused: false}

	timer.PauseToggle()
	assert.True(t, timer.IsPaused)

	timer.PauseToggle()
	assert.False(t, timer.IsPaused)
}

func TestTimer_reset(t *testing.T) {
	timer := &Timer{
		IsPaused:    true,
		IsActive:    true,
		CurrentTodo: "My Todo",
	}

	timer.Reset()

	assert.False(t, timer.IsPaused)
	assert.False(t, timer.IsActive)
	assert.Equal(t, "", timer.CurrentTodo)
}

func TestTimer_NewTimer(t *testing.T) {
	timer := NewTimer(10)

	assert.Equal(t, 10, timer.remaining, "NewTimer should use passed remaining seconds")
}
