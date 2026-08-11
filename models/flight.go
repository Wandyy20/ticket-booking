package models

import "time"

type Flight struct {
	ID            string
	Origin        string
	Destination   string
	DepartureTime time.Time
	ArrivalTime   time.Time
	Price         float64
	TotalSeats    int
}
