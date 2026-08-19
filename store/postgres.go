package store

import (
	"database/sql"
	"errors"
	"fmt"
	"ticket-booking/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) GetAllFlights() ([]models.Flight, error) {
	rows, err := s.db.Query("SELECT id, airline, origin, destination, departure_time, arrival_time, price, total_seats FROM flights")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var f models.Flight
		err := rows.Scan(&f.ID, &f.Airline, &f.Origin, &f.Destination, &f.DepartureTime, &f.ArrivalTime, &f.Price, &f.TotalSeats)

		if err != nil {
			return nil, err
		}
		flights = append(flights, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return flights, nil
}

func (s *PostgresStore) GetFlightByID(id string) (models.Flight, error) {
	var f models.Flight
	err := s.db.QueryRow("SELECT id, airline, origin, destination, departure_time, arrival_time, price, total_seats FROM flights WHERE id = $1", id).Scan(&f.ID, &f.Airline, &f.Origin, &f.Destination, &f.DepartureTime, &f.ArrivalTime, &f.Price, &f.TotalSeats)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Flight{}, errors.New("flight not found")
		}
		return models.Flight{}, err
	}
	return f, nil
}

func (s *PostgresStore) GetSeatsByFlightID(flightID string) ([]models.Seat, error) {
	rows, err := s.db.Query("SELECT id, flight_id, seat_number, is_booked FROM seats WHERE flight_id = $1", flightID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var ss models.Seat
		err := rows.Scan(&ss.ID, &ss.FlightID, &ss.SeatNumber, &ss.IsBooked)
		if err != nil {
			return nil, err
		}
		seats = append(seats, ss)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return seats, nil
}

func (s *PostgresStore) GetSeatByID(seatID string) (models.Seat, error){
	var ss models.Seat
	err := s.db.QueryRow("SELECT id, flight_id, seat_number, is_booked FROM seats WHERE id = $1", seatID).Scan(&ss.ID, &ss.FlightID, &ss.SeatNumber, &ss.IsBooked)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Seat{}, errors.New("No seat found")
		}
		return models.Seat{}, err
	}
	return ss, nil
}

func (s *PostgresStore) CreateBooking(booking models.Booking, seats []models.BookingSeat) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, seat := range seats {
		var isBooked bool

		err := tx.QueryRow("SELECT is_booked FROM seats WHERE id = $1 FOR UPDATE", seat.SeatID).Scan(&isBooked)

		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("seat %s not found", seat.SeatID)
			}
			return err
		}

		if isBooked {
			return fmt.Errorf("seat %s is already booked", seat.SeatID)
		}
	}

	_, err = tx.Exec(
		`INSERT INTO bookings (id, flight_id, status, created_at)
			VALUES ($1, $2, $3, $4)`,
			booking.ID,
			booking.FlightID,
			booking.Status,
			booking.CreatedAt,
	)

	if err != nil {
		return err
	}

	for _, seat := range seats {
		_, err = tx.Exec(`INSERT INTO booking_seats (id, booking_id, seat_id, passenger_name)VALUES ($1, $2, $3, $4)`, seat.ID, seat.BookingID, seat.SeatID, seat.PassengerName)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`UPDATE seats SET is_booked = true WHERE id = $1`, seat.SeatID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresStore) GetBookingByID(id string) (models.Booking, error) {
	var b models.Booking
	err := s.db.QueryRow(`SELECT id, flight_id, status, created_at FROM bookings WHERE id = $1`, id).Scan(&b.ID, &b.FlightID, &b.Status, &b.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Booking{}, errors.New("no booking found")
		}
		return models.Booking{}, err
	}
	return b, err
}

func (s *PostgresStore) CancelBooking(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var status string 

	err = tx.QueryRow(`SELECT status FROM bookings WHERE id = $1 FOR UPDATE`, id,).Scan(&status)

	if err != nil {
		return err
	}

	if status == "cancelled" {
		return errors.New("booking already cancelled")
	}

	_, err = tx.Exec(`UPDATE bookings SET status = $1 WHERE id = $2`, "cancelled", id)

	if err != nil {
		return err
	}

	rows, err := tx.Query(`SELECT seat_id FROM booking_seats WHERE booking_id = $1`, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("booking not found")
		}
		return err
	}

	defer rows.Close()

	for rows.Next( ){
		var seatID string
		err := rows.Scan(&seatID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`UPDATE seats SET is_booked = false WHERE id = $1`, seatID)

		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return tx.Commit()
}