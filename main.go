package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-books-api/db"
	"go-books-api/routes/auth"
	"go-books-api/routes/books"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello Go Books API")
	mysql := db.InitMySQL()
	errMySQL := mysql.Ping()
	if errMySQL != nil {
		fmt.Printf("Ping didn't worked %s", errMySQL)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		log.Printf("Welcome")
	})

	auth.AuthRoutes(mysql, r)
	books.BooksRoutes(mysql, r)
	errServer := http.ListenAndServe(":80", r)
	if errServer != nil {
		log.Printf("Server Error %s", errServer)
	}
}
