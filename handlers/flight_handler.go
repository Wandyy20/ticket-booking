package handlers

import (
	"encoding/json"
	"net/http"
	"ticket-booking/store"
	"github.com/go-chi/chi/v5"
)

func ListFlightHandler(s store.FlightStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flights, err := s.GetAllFlights()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flights)
	}
}

func GetFlightByIDHandler(s store.FlightStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flightID := chi.URLParam(r, "id")
		flights, err := s.GetFlightByID(flightID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flights)
	}
}

func GetFlightSeatsHandler(s store.SeatStore) http.HandlerFunc {
	return func( w http.ResponseWriter, r *http.Request){
		flightID := chi.URLParam(r, "id")
		seats, err := s.GetSeatByFlightID(flightID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(seats)
	}
}

type BookingSeatRequest struct {
	SeatID string `json:"seatID"`
	PassengerName string `json:"passengerName"`
}

type BookingRequest struct {
	FlightID string `json:"flightID"`
	Seats []BookingSeatRequest `json:"seats"`
}


