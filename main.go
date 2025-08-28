package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-books-api/db"
	"go-books-api/routes/auth"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello Go Books API")
	mysql := db.InitMySQL()
	err := mysql.Ping()
	if err != nil {
		fmt.Printf("Ping didn't worked %s", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		log.Printf("Welcome")
	})

	auth.AuthRoutes(mysql, r)
	http.ListenAndServe(":80", r)
}
