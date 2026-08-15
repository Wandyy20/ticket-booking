package handlers

import (
	"encoding/json"
	"net/http"
	"ticket-booking/models"
	"ticket-booking/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateBookingHandler(s store.BookingStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var req BookingRequest

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		newBooking := models.Booking{
			ID: uuid.NewString(),
			FlightID: req.FlightID,
			Status: models.StatusConfirmed,
			CreatedAt: time.Now(),
		}

		var newSeats []models.BookingSeat
		for _, seatReq := range req.Seats {
			newSeats = append(newSeats, models.BookingSeat{
				ID: uuid.NewString(),
				BookingID: newBooking.ID,
				SeatID: seatReq.SeatID,
				PassengerName: seatReq.PassengerName,
			})
		}
		err = s.CreateBooking(newBooking, newSeats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBooking)
	}
}

	func GetBookingHandler(s store.BookingStore) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			bookID := chi.URLParam(r, "id")
			book, err := s.GetBookingByID(bookID)

			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
		}
	}

	func DeleteBookingHandler(s store.BookingStore)http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			bookID := chi.URLParam(r, "id")
			err := s.CancelBooking(bookID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("booking cancelled successfully"))
		}
	}