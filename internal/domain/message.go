package domain

import "time"

type Message struct {
	ID        string
	Content   string
	Timestamp time.Time
}
