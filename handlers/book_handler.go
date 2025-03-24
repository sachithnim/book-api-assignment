package handlers

import (
	"book-api-assignment/models"
	"book-api-assignment/repository"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// Standard response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Utility function to send JSON response
func sendResponse(w http.ResponseWriter, statusCode int, status string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json") // Ensure correct response type
	w.WriteHeader(statusCode)

	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	// Send JSON response
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.LoadBooks()
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "error", "Failed to load books", nil)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))   // Get page number (pagination)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Get limit per page

	if page < 1 {
		page = 1 // Default to page 1
	}
	if limit < 1 {
		limit = 5 // Default limit per page
	}

	start := (page - 1) * limit
	end := start + limit
	totalBooks := len(books)

	if start >= totalBooks {
		sendResponse(w, http.StatusOK, "success", "No books found for this page", []models.Book{})
		return
	}

	if end > totalBooks {
		end = totalBooks
	}

	// Paginated books
	paginatedBooks := books[start:end]

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalBooks) / float64(limit)))

	// Create pagination response
	paginationResponse := map[string]interface{}{
		"status":      "success",
		"message":     "Books retrieved successfully",
		"totalBooks":  totalBooks,
		"totalPages":  totalPages,
		"currentPage": page,
		"limit":       limit,
		"data":        paginatedBooks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginationResponse)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := repository.FindBookByID(params["id"])
	if err != nil {
		sendResponse(w, http.StatusNotFound, "error", "Book not found", nil)
		return
	}
	sendResponse(w, http.StatusOK, "success", "Book retrieved successfully", book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Ensure BookID is provided
	if book.BookID == "" {
		sendResponse(w, http.StatusBadRequest, "error", "BookID is required", nil)
		return
	}

	books, _ := repository.LoadBooks()

	// Check for duplicate BookID
	for _, existingBook := range books {
		if existingBook.BookID == book.BookID {
			sendResponse(w, http.StatusConflict, "error", "BookID already exists", nil)
			return
		}
	}

	books = append(books, book)
	repository.SaveBooks(books)

	sendResponse(w, http.StatusCreated, "success", "Book created successfully", book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedBook models.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	books, _ := repository.LoadBooks()
	for i, book := range books {
		if book.BookID == params["id"] {
			// Only update non empty fields
			if updatedBook.Title != "" {
				book.Title = updatedBook.Title
			}
			if updatedBook.AuthorID != "" {
				book.AuthorID = updatedBook.AuthorID
			}
			if updatedBook.PublisherID != "" {
				book.PublisherID = updatedBook.PublisherID
			}
			if updatedBook.PublicationDate != "" {
				book.PublicationDate = updatedBook.PublicationDate
			}
			if updatedBook.ISBN != "" {
				book.ISBN = updatedBook.ISBN
			}
			if updatedBook.Pages != 0 {
				book.Pages = updatedBook.Pages
			}
			if updatedBook.Genre != "" {
				book.Genre = updatedBook.Genre
			}
			if updatedBook.Description != "" {
				book.Description = updatedBook.Description
			}
			if updatedBook.Price != 0 {
				book.Price = updatedBook.Price
			}
			if updatedBook.Quantity != 0 {
				book.Quantity = updatedBook.Quantity
			}

			books[i] = book
			repository.SaveBooks(books)

			sendResponse(w, http.StatusOK, "success", "Book updated successfully", book)
			return
		}
	}

	sendResponse(w, http.StatusNotFound, "error", "Book not found", nil)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	books, _ := repository.LoadBooks()
	index := -1
	for i, book := range books {
		if book.BookID == params["id"] {
			index = i
			break
		}
	}

	if index == -1 {
		sendResponse(w, http.StatusNotFound, "error", "Book not found", nil)
		return
	}

	// Remove book
	books = append(books[:index], books[index+1:]...)

	// Save changes
	if err := repository.SaveBooks(books); err != nil {
		sendResponse(w, http.StatusInternalServerError, "error", "Failed to delete book", nil)
		return
	}

	sendResponse(w, http.StatusOK, "success", "Book deleted successfully", nil)
}

// Search for books by title and description (concurrent implementation)
func SearchBooksHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q"))) // Get query parameter and convert to lowercase
	if query == "" {
		sendResponse(w, http.StatusBadRequest, "error", "Query parameter 'q' is required", nil)
		return
	}

	books, err := repository.LoadBooks()
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "error", "Failed to load books", nil)
		return
	}

	// Number of goroutines to use for searching
	numWorkers := 4
	chunkSize := (len(books) + numWorkers - 1) / numWorkers

	// Channel to collect results
	resultsChan := make(chan []models.Book, numWorkers)
	var wg sync.WaitGroup

	// Launch concurrent search workers
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(books) {
			end = len(books)
		}

		wg.Add(1)
		go searchBooksWorker(books[start:end], query, resultsChan, &wg)
	}

	// Wait for all workers to complete
	wg.Wait()
	close(resultsChan)

	// Merge results
	var matchedBooks []models.Book
	for books := range resultsChan {
		matchedBooks = append(matchedBooks, books...)
	}

	sendResponse(w, http.StatusOK, "success", "Search completed", matchedBooks)
}

// Worker function to search a subset of books
func searchBooksWorker(books []models.Book, query string, resultsChan chan<- []models.Book, wg *sync.WaitGroup) {
	defer wg.Done()

	var matchedBooks []models.Book
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), query) || strings.Contains(strings.ToLower(book.Description), query) {
			matchedBooks = append(matchedBooks, book)
		}
	}

	// Send results to the channel
	resultsChan <- matchedBooks
}
