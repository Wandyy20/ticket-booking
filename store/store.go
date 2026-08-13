package store

import "ticket-booking/models"

type FlightStore interface {
	GetAllFlights() ([]models.Flight, error)
	GetFlightByID(id string) (models.Flight, error)
}

type SeatStore interface {
	GetSeatByFlightID(flightID string) ([]models.Seat, error)
	GetSeatByID(id string) (models.Seat, error)
}

type BookingStore interface {
	CreateBooking(booking models.Booking, seats []models.BookingSeat) error
	GetBookingByID(id string) (models.Booking, error)
	CancelBooking(id string) error
}