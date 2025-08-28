package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("Hashing failed %s", err)
	}
	return string(bytes)
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
