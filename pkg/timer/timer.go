package timer

import "time"

type Timer interface {
	Now() time.Time
}

type UTCTimer struct{}

func NewUTCTimer() *UTCTimer {
	return &UTCTimer{}
}

func (s UTCTimer) Now() time.Time {
	return time.Now().UTC()
}
