package repository

import (
	"book-api-assignment/models"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const dataFile = "data/books.json"

var mu sync.Mutex

// reads books from the JSON file.
func LoadBooks() ([]models.Book, error) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(dataFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.Book{}, nil // If file does not exist, return empty
		}
		return nil, err // Return actual error
	}

	var books []models.Book
	err = json.Unmarshal(file, &books)
	if err != nil {
		return nil, err // Return error if JSON parsing fails
	}
	return books, nil
}

// writes books to the JSON file.
func SaveBooks(books []models.Book) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

// Find a book by ID
func FindBookByID(id string) (*models.Book, error) {
	books, _ := LoadBooks()
	for _, book := range books {
		if book.BookID == id {
			return &book, nil
		}
	}
	return nil, errors.New("book not found")
}
