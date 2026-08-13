package store

import (
	"fmt"
	"time"
	"strings"

	"ticket-booking/models"
)

func (s *InMemoryStore) Seed() {

	flights := []models.Flight{
		{
			ID:            "FL001",
			Airline:       "Lion Air",
			Origin:        "Jakarta",
			Destination:   "Bali",
			DepartureTime: time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC),
			Price:         800000,
			TotalSeats:    20,
		},
		{
			ID:            "FL002",
			Airline:       "Garuda Indonesia",
			Origin:        "Jakarta",
			Destination:   "Surabaya",
			DepartureTime: time.Date(2026, 9, 2, 9, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 9, 2, 10, 30, 0, 0, time.UTC),
			Price:         950000,
			TotalSeats:    20,
		},
		{
			ID:            "FL003",
			Airline:       "Citilink",
			Origin:        "Jakarta",
			Destination:   "Yogyakarta",
			DepartureTime: time.Date(2026, 9, 3, 7, 30, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 9, 3, 9, 0, 0, 0, time.UTC),
			Price:         700000,
			TotalSeats:    20,
		},
		{
			ID:            "FL004",
			Airline:       "AirAsia",
			Origin:        "Bandung",
			Destination:   "Bali",
			DepartureTime: time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 9, 4, 12, 30, 0, 0, time.UTC),
			Price:         850000,
			TotalSeats:    20,
		},
		{
			ID:            "FL005",
			Airline:       "Batik Air",
			Origin:        "Jakarta",
			Destination:   "Medan",
			DepartureTime: time.Date(2026, 9, 5, 13, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 9, 5, 15, 30, 0, 0, time.UTC),
			Price:         1200000,
			TotalSeats:    20,
		},
	}

	for _, flight := range flights {

		s.flights[flight.ID] = flight

		for row := 1; row <= 10; row++ {
			for _, column := range []string{"A", "B"} {

				seatNumber := fmt.Sprintf("%d%s", row, column)

				seatID := fmt.Sprintf(
					"%s-seat-%s",
					flight.ID,
					strings.ToLower(seatNumber),
				)

				seat := models.Seat{
					ID:         seatID,
					FlightID:   flight.ID,
					SeatNumber: seatNumber,
					IsBooked:   false,
				}

				s.seats[seat.ID] = seat
			}
		}
	}
}