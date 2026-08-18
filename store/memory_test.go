package store

import (
	"sync"
	"testing"
	"ticket-booking/models"
	"fmt"
)

func TestGetAllFlights(t *testing.T) {

	s := NewInMemoryStore()

	s.Seed()

	flights, err := s.GetAllFlights()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(flights) == 0 {
		t.Errorf("expected flights, got empty list")
	}
}

func TestCreateBooking_Success (t *testing.T) {
	s := NewInMemoryStore()
	s.Seed()

	booking := models.Booking{
		ID:       "book-test-1",
		FlightID: "FL001",
		Status:   models.StatusConfirmed,
	}
	seats := []models.BookingSeat{
		{ID: "bs-test-1", BookingID: "book-test-1", SeatID: "FL001-seat-1a", PassengerName: "Budi"},
	}

	err := s.CreateBooking(booking, seats)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	seat, _ := s.GetSeatByID("FL001-seat-1a")
	if !seat.IsBooked {
		t.Errorf("expected seat to be booked, but it's not")
	}
}

func TestCreateBooking_SeatAlreadyBooked(t *testing.T) {
	s := NewInMemoryStore()
	s.Seed()

	s.CreateBooking(models.Booking{ID: "book-1", FlightID: "FL001"}, []models.BookingSeat{
		{ID: "bs-1", BookingID: "book-1", SeatID: "FL001-seat-1a", PassengerName: "Budi"},
	})

	err := s.CreateBooking(models.Booking{ID: "book-2", FlightID: "FL001"}, []models.BookingSeat{
		{ID: "bs-2", BookingID: "book-2", SeatID: "FL001-seat-1a", PassengerName: "Ani"},
	})

	if err == nil {
		t.Errorf("expected error (seat already booked), got nil")
	}
}

func TestCreateBooking_RaceCondition(t *testing.T) {
	s := NewInMemoryStore()
	s.Seed()

	numAttempts := 10  
	var wg sync.WaitGroup
	var mu sync.Mutex    
	successCount := 0

	for i := 0; i < numAttempts; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			booking := models.Booking{
				ID:       fmt.Sprintf("book-%d", id),
				FlightID: "FL001",
				Status:   models.StatusConfirmed,
			}
			seats := []models.BookingSeat{
				{
					ID:            fmt.Sprintf("bs-%d", id),
					BookingID:     booking.ID,
					SeatID:        "FL001-seat-1a",  
					PassengerName: fmt.Sprintf("User-%d", id),
				},
			}

			err := s.CreateBooking(booking, seats)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()  

	if successCount != 1 {
		t.Errorf("expected exactly 1 successful booking, got %d", successCount)
	}
}