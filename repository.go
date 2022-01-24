package libapi

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

//NewRepository connects to the sql db
func NewRepository(db *sql.DB) Repository {
	return Repository{
		db: db,
	}
}

type Repository struct {
	db *sql.DB
}

//Creates a book in the library
func (r Repository) Create(book Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO book (id, name, author, isbn) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id = ?, name = ?, author = ?, isbn = ?`
	log.Info("book in repo: ", book)
	_, err := r.db.ExecContext(
		ctx,
		query,
		book.ID,
		book.Name,
		book.Author,
		book.ISBN,
		book.ID, // start of upsert
		book.Name,
		book.Author,
		book.ISBN,
	)
	if err != nil {
		log.Error("Error while storing model")
		return err
	}

	return nil
}

// FindByID attaches the book repository and find data based on id
func (r Repository) FindByID(id string) (*Book, error) {
	book := new(Book)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, "SELECT id, name, author, isbn FROM book WHERE id = ?", id).Scan(&book.ID, &book.Name, &book.Author, &book.ISBN)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// Find attaches the book repository and find all data
func (r Repository) Find() ([]Book, error) {
	books := make([]Book, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, "SELECT id, name, author, isbn FROM book")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		book := new(Book)
		err = rows.Scan(
			&book.ID,
			&book.Name,
			&book.Author,
			&book.ISBN,
		)

		if err != nil {
			return nil, err
		}
		books = append(books, *book)
	}

	return books, nil
}

// Update attaches the book repository and update data based on id
func (r Repository) Update(book Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE book SET name = ?, author = ?, isbn = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, book.Name, book.Author, book.ISBN, book.ID)
	return err
}

// Delete attaches the book repository and delete data based on id
func (r Repository) Delete(id string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM books WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}
