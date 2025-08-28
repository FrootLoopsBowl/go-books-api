package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitMySQL() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("mysql", os.Getenv("MYSQL_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Error connecting to database %s", err)
	}
	return db
}
