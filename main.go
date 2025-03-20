package main

import (
	"book-api-assignment/routes"
	"log"
	"net/http"
)

func main() {
	router := routes.SetupRoutes()
	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
