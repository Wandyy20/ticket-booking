package store

import (
	"errors"
	"sync"
	"ticket-booking/models"
)

type InMemoryStore struct {
	mu           sync.Mutex
	flightLocks  map[string]*sync.Mutex
	flights      map[string]models.Flight
	seats        map[string]models.Seat
	bookings     map[string]models.Booking
	bookingSeats map[string]models.BookingSeat
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		flights:      make(map[string]models.Flight),
		flightLocks:  make(map[string]*sync.Mutex),
		seats:        make(map[string]models.Seat),
		bookings:     make(map[string]models.Booking),
		bookingSeats: make(map[string]models.BookingSeat),
	}
}

func (s *InMemoryStore) getFlightLock(flightID string) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.flightLocks[flightID]; !ok {
		s.flightLocks[flightID] = &sync.Mutex{}
	}
	return s.flightLocks[flightID]
}

func (s *InMemoryStore) GetAllFlights() ([]models.Flight, error){
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []models.Flight{}
	for _, f := range s.flights {
		result = append(result, f)
	}
	return result, nil
}

func (s *InMemoryStore) GetFlightByID(id string) (models.Flight, error){
	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.flights[id]
	if !ok {
		return models.Flight{}, errors.New("no flights found")
	}
	return f, nil
}

func (s *InMemoryStore) GetSeatsByFlightID(flightID string) ([]models.Seat, error){
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []models.Seat{}
	for _, seat := range s.seats {
		if seat.FlightID == flightID {
			result = append(result, seat)
		}
	}
	return result, nil
}

func (s *InMemoryStore) GetSeatByID(id string) (models.Seat, error){
	s.mu.Lock()
	defer s.mu.Unlock()

	f, ok := s.seats[id]
	if !ok {
		return models.Seat{}, errors.New("seat not found")
	}
	return f, nil
}

func (s *InMemoryStore) GetBookingByID(id string) (models.Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.bookings[id]
	if !ok {
		return models.Booking{}, errors.New("booking not found")
	}
	return b, nil
}

func (s *InMemoryStore) CancelBooking(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	booking, ok := s.bookings[id]
	if !ok {
		return errors.New("booking not found")
	}

	booking.Status = models.StatusCancelled
	s.bookings[id] = booking

	for _, bs := range s.bookingSeats{
		if bs.BookingID == id {
			seat := s.seats[bs.SeatID]
			seat.IsBooked = false
			s.seats[bs.SeatID] = seat
		}
	}
	return nil
}

func (s *InMemoryStore) CreateBooking(booking models.Booking, seats []models.BookingSeat) error {
	flightLock := s.getFlightLock(booking.FlightID)
	flightLock.Lock()
	defer flightLock.Unlock()

	for _, bs := range seats {
		seat, ok := s.seats[bs.SeatID]
		if !ok {
			return errors.New("seat not found")
		}
		if seat.IsBooked {
			return  errors.New("seat already booked")
		}
	}

	s.bookings[booking.ID] = booking
	for _, bs := range seats {
		s.bookingSeats[bs.ID] = bs
		seat := s.seats[bs.SeatID]
		seat.IsBooked = true
		s.seats[bs.SeatID] = seat
	}
	return nil	
}
