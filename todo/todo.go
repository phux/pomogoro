package todo

import "time"

type Todo struct {
	title     string
	trackings []Tracking
}

func NewTodo(title string) *Todo {
	trackings := make([]Tracking, 0)
	return &Todo{title: title, trackings: trackings}
}

func (t Todo) Title() string {
	return t.title
}

func (t Todo) TrackedTime() time.Duration {
	var total time.Duration
	for _, tracking := range t.trackings {
		total += tracking.duration
	}

	return total
}

func (t *Todo) Track(duration time.Duration) {
	t.trackings = append(t.trackings, NewTracking(duration))
}

func (t *Todo) TrackWithDate(duration time.Duration, date time.Time) {
	t.trackings = append(t.trackings, NewTrackingWithDate(duration, date))
}
