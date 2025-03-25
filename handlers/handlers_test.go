package handlers

import (
	"book-api-assignment/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateBook(t *testing.T) {
	// Prepare test data
	book := models.Book{
		Title:           "Test Book",
		AuthorID:        "test-author-id",
		PublisherID:     "test-publisher-id",
		PublicationDate: "2025-03-25",
		ISBN:            "1234567890",
		Pages:           200,
		Genre:           "Fiction",
		Description:     "A test book description.",
		Price:           20.99,
		Quantity:        10,
	}

	// Marshal the book data into JSON
	jsonBook, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Error marshaling book: %v", err)
	}

	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBook)) // Create a new request to the CreateBook endpoint
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder() // Create a ResponseRecorder to capture the response

	router := mux.NewRouter() // Set up the router and handle the request
	router.HandleFunc("/books", CreateBook).Methods("POST")

	// Serve the request and check the response
	router.ServeHTTP(rr, req)

	// Assert that the status code is what we expect
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, rr.Code)
	}

	// Assert that the response body contains the correct information
	var response Response
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Check response status is 'success'
	if response.Status != "success" {
		t.Errorf("Expected status 'success', but got '%s'", response.Status)
	}

	// Verify the created book's details in the response
	if response.Data == nil {
		t.Error("Expected book data in the response, but got nil")
	} else {
		var createdBook models.Book
		data, _ := json.Marshal(response.Data)
		err := json.Unmarshal(data, &createdBook)
		if err != nil {
			t.Fatalf("Error unmarshalling book data: %v", err)
		}

		if createdBook.Title != book.Title {
			t.Errorf("Expected book title '%s', but got '%s'", book.Title, createdBook.Title)
		}
	}
}
