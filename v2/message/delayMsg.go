package message

import "time"

type DelayMessage struct {
	msg Message
	t time.Time
}
