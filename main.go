package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ticket-booking/handlers"
	"ticket-booking/store"
)
func main(){
	appStore := store.NewInMemoryStore()
	appStore.Seed()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/flights", handlers.ListFlightHandler(appStore))
	r.Get("/flights/{id}", handlers.GetFlightByIDHandler(appStore))
	r.Get("/flights/{id}/seats", handlers.GetFlightSeatsHandler(appStore))

	r.Post("/bookings", handlers.CreateBookingHandler(appStore))
	r.Get("/bookings/{id}", handlers.GetBookingHandler(appStore))
	r.Delete("/bookings/{id}", handlers.DeleteBookingHandler(appStore))

	log.Println("Server jalan di :8080")
	http.ListenAndServe(":8080", r)
}