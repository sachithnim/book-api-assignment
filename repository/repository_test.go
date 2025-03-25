package repository

import (
	"book-api-assignment/models"
	"encoding/json"
	"os"
	"testing"
)

func TestLoadBooks(t *testing.T) {
	// Create the data directory if it doesn't exist
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating data directory: %v", err)
	}

	// Add test data to the books.json
	testBooks := []models.Book{
		{
			BookID:          "fab04289-df6c-4ae7-ace4-b867edfc18d3",
			Title:           "Test Book",
			AuthorID:        "author-id",
			PublisherID:     "publisher-id",
			PublicationDate: "2025-03-25",
			ISBN:            "1234567890",
			Pages:           200,
			Genre:           "Fiction",
			Description:     "A test book",
			Price:           20.99,
			Quantity:        10,
		},
	}

	// Save test data to the JSON file
	data, err := json.Marshal(testBooks)
	if err != nil {
		t.Fatalf("Error marshalling test books: %v", err)
	}

	err = os.WriteFile("data/books.json", data, 0644)
	if err != nil {
		t.Fatalf("Error writing to books.json: %v", err)
	}

	// Call LoadBooks function
	books, err := LoadBooks()
	if err != nil {
		t.Fatalf("Error loading books: %v", err)
	}

	// Check if books are returned
	if len(books) == 0 {
		t.Error("Expected some books, got none")
	}

	// Optionally clean up after test
	err = os.Remove("data/books.json")
	if err != nil {
		t.Fatalf("Error removing test books file: %v", err)
	}

	err = os.RemoveAll("data") // Remove the entire directory
	if err != nil {
		t.Fatalf("Error removing test data directory: %v", err)
	}
}
