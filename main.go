package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	// ===================================
	// connect to database
	// ===================================

	// connect to the dbs
	mySQLDB := connectToDB()
	if migErr := handleMigrations(mySQLDB); migErr != nil {
		log.Error("migration error: ", migErr)
	}

	repo := NewRepository(mySQLDB)

	// initialize a panic watcher
	defer recover()

	// Init router
	r := mux.NewRouter()

	//decided to write middleware to handle setting header
	//this didn't seem sexy, only functional after reading docs
	//https://github.com/gorilla/mux/blob/v1.8.0/middleware.go#L11
	//https://pkg.go.dev/github.com/gorilla/mux#section-readme

	r.Use(setHeaderContentMiddleware)

	r.HandleFunc("/books", getBooks(repo)).Methods("GET")
	r.HandleFunc("/books/{id}", getBook(repo)).Methods("GET")
	r.HandleFunc("/books", createBook(repo)).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook(repo)).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook(repo)).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))

}

func setHeaderContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
