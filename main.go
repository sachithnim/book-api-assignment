package main

import (
	"book-api-assignment/handlers"
	"book-api-assignment/routes"
	"log"
	"net/http"
)

func main() {
	router := routes.SetupRoutes()

	// Ensure this line is added to handle the root route.
	router.HandleFunc("/", handlers.RootHandler).Methods("GET")

	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
