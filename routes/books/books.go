package books

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"go-books-api/models"
	"log"
	"net/http"
)

func BooksRoutes(db *sql.DB, r *mux.Router) {
	bRouter := r.PathPrefix("/books").Subrouter()

	bRouter.HandleFunc("/create/{name}/{author}/{category}/{image}", func(w http.ResponseWriter, r *http.Request) { //TODO: FIX IMAGE CAN'T USE LINK | COULD REMOVE AND USE SERVER STORAGE
		vars := mux.Vars(r)
		books := models.Books{
			Name: vars["name"], Author: vars["author"], Image: vars["image"], Category: vars["category"],
		}
		query := "SELECT name, author FROM books WHERE name = ? AND author = ?"
		var bName, bAuthor string
		err := db.QueryRow(query, vars["name"], vars["author"]).Scan(&bName, &bAuthor)
		if errors.Is(err, sql.ErrNoRows) {
			log.Print("Creating book")
			_, err := db.Exec("INSERT INTO books (name, author, image, category) VALUES (?, ?, ?, ?)", books.Name, books.Author, books.Image, books.Category)
			if err != nil {
				log.Printf("[CREATE BOOKS] Query error %s", err)
			} else {
				log.Print("Book created")
				w.WriteHeader(http.StatusCreated)
			}
		} else if err != nil {
			log.Printf("[CREATE BOOKS] Query error %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Print("Book already existed")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}).Methods("POST")

	bRouter.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		author := r.URL.Query().Get("author")
		category := r.URL.Query().Get("category")

		query := "SELECT * FROM books WHERE 1=1"
		if name != "" {
			query += " AND name = " + "'" + name + "'"
		}
		if author != "" {
			query += " AND author = " + "'" + author + "'"
		}
		if category != "" {
			query += " AND category = " + "'" + category + "'"
		}

		rows, err := db.Query(query)
		defer rows.Close()

		var books []models.Books
		for rows.Next() {
			var b models.Books
			err := rows.Scan(&b.Id, &b.Name, &b.Author, &b.Category, &b.Image)
			if err != nil {
				log.Printf("[GET BOOKS] Query error %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			books = append(books, b)
		}

		err = rows.Err()
		if err != nil {
			log.Printf("[GET BOOKS] Query error %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		if len(books) > 0 {
			for i := 0; i < len(books); i++ {
				log.Printf("%+v", books[i])
			}
		} else {
			log.Print("Found no book")
			w.WriteHeader(http.StatusNotFound)
		}
	}).Methods("GET")
}
