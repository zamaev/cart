package model

import "time"

type Event struct {
	OrderID int64
	Status  string
	Time    time.Time
}
