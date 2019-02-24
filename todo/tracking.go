package todo

import "time"

type Tracking struct {
	duration time.Duration
	date     time.Time
}

func NewTracking(duration time.Duration) Tracking {
	return Tracking{
		duration: duration,
	}
}

func NewTrackingWithDate(duration time.Duration, date time.Time) Tracking {
	tracking := NewTracking(duration)
	tracking.date = date

	return tracking
}
