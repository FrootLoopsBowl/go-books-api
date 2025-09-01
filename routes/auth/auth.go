package auth

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"go-books-api/models"
	"log"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
    	"github.com/joho/godotenv"
	"os"
	"encoding/json"
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
				tokenString := CreateJWT(username)

				response := map[string]string{"token": tokenString}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)

				log.Printf("Welcome %s", username)
			} else {
				log.Printf("Login Denied")
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	}).Methods("POST")

	aRouter.HandleFunc("/testtoken/{token}", func(w http.ResponseWriter, r *http.Request) { //Debug method
		vars := mux.Vars(r)
		token := vars["token"]
		result := VerifyToken(token)
		if result {
			log.Print("Token is good!!")
		} else {
			log.Print("Token is no good")
		}
	})
}



func CreateJWT(username string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Error with the JWT generation %s", err)
	}
	return tokenString
}

func VerifyToken(tokenString string) bool {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		log.Printf("Error with JWT %s", err)
	}
	return token.Valid
}
