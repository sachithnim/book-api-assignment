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

func LoadBooks() ([]models.Book, error) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(dataFile)
	if err != nil {
		return []models.Book{}, nil
	}

	var books []models.Book
	err = json.Unmarshal(file, &books)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func SaveBooks(books []models.Book) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func FindBookByID(id string) (*models.Book, error) {
	books, _ := LoadBooks()
	for _, book := range books {
		if book.BookID == id {
			return &book, nil
		}
	}
	return nil, errors.New("book not found")
}
