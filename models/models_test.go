package models

import "testing"

func TestBookModel(t *testing.T) {
	book := Book{
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

	if book.Title != "Test Book" {
		t.Errorf("Expected 'Test Book', got '%s'", book.Title)
	}
}
