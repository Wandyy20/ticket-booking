package models

import "time"

type BookingStatus string

const (
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID        string
	FlightID  string
	Status    BookingStatus
	CreatedAt time.Time
}