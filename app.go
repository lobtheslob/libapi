package libapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Init var books as a slice of Book struct
var books []Book

// Get all books
func getBooks(repo Repository) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		books, err := repo.Find()
		if err != nil {
			json.NewEncoder(w).Encode(books)
		}
	}
}

// Get single book
func getBook(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r) // Gets params
		for _, item := range books {
			itemID, _ := strconv.ParseInt(params["id"], 10, 32)
			if string(rune(itemID)) == params["id"] {
				item, err := repo.FindByID(string(rune(item.ID)))
				if err != nil { 
					json.NewEncoder(w).Encode(item)
					return
				}
			}
		}
		json.NewEncoder(w).Encode(&Book{})
	}
}

// Add new book
func createBook(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var book Book
		_ = json.NewDecoder(r.Body).Decode(&book)
		err := repo.Create(book)
		if err != nil {
		  return 
		}
		books = append(books, book)
		json.NewEncoder(w).Encode(book)
	}
}

// Delete book
func deleteBook(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for idx, item := range books {
			bookID, _ := strconv.ParseInt(params["id"], 10, 32)
			item.ID = int(bookID)
			if string(rune(item.ID)) == params["id"] {
				err :=  repo.Delete(string(rune(item.ID)))
				if err != nil {
					break
				}
				books = append(books[:idx], books[idx+1:]...)
				break
			}
		}
		json.NewEncoder(w).Encode(books)
	}
}
// Update book
func updateBook(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for idx, item := range books {
				bookID, _ := strconv.ParseInt(params["id"], 10, 32)
				item.ID = int(bookID)
			if string(rune(item.ID)) == params["id"] {
				books = append(books[:idx], books[idx+1:]...)
				var book Book
				_ = json.NewDecoder(r.Body).Decode(&book)
				err := repo.Update(book)
				if err != nil {
					return
				} 
				books = append(books, book)
				json.NewEncoder(w).Encode(book)
				return
			}
		}
	}
}