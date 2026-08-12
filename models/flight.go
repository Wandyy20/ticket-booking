package models

import "time"

type Flight struct {
	ID            string
	Airline       string
	FlightNumber  string
	Origin        string
	Destination   string
	DepartureTime time.Time
	ArrivalTime   time.Time
	Price         float64
	TotalSeats    int
}
