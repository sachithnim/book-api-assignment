package handlers

import (
	"book-api-assignment/models"
	"book-api-assignment/repository"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books, _ := repository.LoadBooks()
	json.NewEncoder(w).Encode(books)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := repository.FindBookByID(params["id"])
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	// Ensure BookID is provided
	if book.BookID == "" {
		http.Error(w, "BookID is required", http.StatusBadRequest)
		return
	}

	books, _ := repository.LoadBooks()

	// Check for duplicate BookID
	for _, existingBook := range books {
		if existingBook.BookID == book.BookID {
			http.Error(w, "BookID already exists", http.StatusConflict)
			return
		}
	}

	books = append(books, book)
	repository.SaveBooks(books)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedBook models.Book
	json.NewDecoder(r.Body).Decode(&updatedBook)

	books, _ := repository.LoadBooks()
	for i, book := range books {
		if book.BookID == params["id"] {
			books[i] = updatedBook
			repository.SaveBooks(books)
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	books, _ := repository.LoadBooks()
	for i, book := range books {
		if book.BookID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			repository.SaveBooks(books)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}
