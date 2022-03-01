package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Init var books as a slice of Book struct
var books []Book

// Get all books
func getBooks(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		books, err := repo.Find()
		if err == nil {
			json.NewEncoder(w).Encode(books)
		}
	}
}

// Get single book
func getBook(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r) // Gets params
		books, err := repo.Find()
		log.Info(err)
		for _, item := range books {
			itemID, _ := strconv.ParseInt(params["id"], 10, 32)
			log.Info("inside get book, itemID: ", itemID)
			if int64(item.ID) == itemID {
				item, err := repo.FindByID(params["id"])
				if err == nil {
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
		books, _ := repo.Find()
		for idx, item := range books {
			i, _ := strconv.ParseInt(params["id"], 10, 64)
			if int64(item.ID) == i {
				err := repo.Delete(params["id"])
				log.Error(err)
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
		books, _ := repo.Find()
		for idx, item := range books {
			i, _ := strconv.ParseInt(params["id"], 10, 64)
			if int64(item.ID) == i {
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
