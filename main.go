package main

import (
	"book-api-assignment/handlers"
	"book-api-assignment/middleware"
	"book-api-assignment/routes"
	"log"
	"net/http"
)

func main() {
	router := routes.SetupRoutes()

	// Add the logging middleware
	router.Use(middleware.Logger)

	// root route.
	router.HandleFunc("/", handlers.RootHandler).Methods("GET")

	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
