package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ticket-booking/store"
)
func main(){
	appStore := store.NewInMemoryStore()
	appStore.Seed()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, req *http.Request){
		w.Write([]byte("pong"))
	})

	log.Println("Server jalan di :8080")
	http.ListenAndServe(":8080", r)
}