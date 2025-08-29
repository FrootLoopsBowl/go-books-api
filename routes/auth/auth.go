package auth

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"go-books-api/models"
	"log"
	"net/http"
	"time"
)

func AuthRoutes(db *sql.DB, r *mux.Router) {
	aRouter := r.PathPrefix("/auth").Subrouter()

	aRouter.HandleFunc("/create/{username}/{password}", func(w http.ResponseWriter, r *http.Request) { //TODO: JSON SIGNUP
		vars := mux.Vars(r)
		password := vars["password"]
		passwordHash := HashPassword(password)
		user := models.User{
			Username: vars["username"], Password: passwordHash, CreatedAt: time.Now(),
		}
		_, err := db.Exec("INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)", user.Username, user.Password, user.CreatedAt)
		if err != nil {
			log.Printf("User creation failed %s", err)
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	aRouter.HandleFunc("/login/{username}/{password}", func(w http.ResponseWriter, r *http.Request) { //TODO: JSON LOGIN
		vars := mux.Vars(r)
		username := vars["username"]
		passwordRaw := vars["password"]
		query := "SELECT password FROM users WHERE username = ?"
		var password string
		err := db.QueryRow(query, username).Scan(&password)
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Didn't find a user")
			w.WriteHeader(http.StatusUnauthorized)
		} else if err != nil {
			log.Printf("Login query failed %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if CheckPassword(passwordRaw, password) {
				log.Printf("Welcome %s", username)
				w.WriteHeader(http.StatusAccepted)
			} else {
				log.Printf("Login Denied")
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	}).Methods("POST")
}
