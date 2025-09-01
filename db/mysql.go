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

func CreateDatabase(db *sql.DB) {
	queryAuth := `
		CREATE TABLE users (
			id INT AUTO_INCREMENT,
			username VARCHAR(80) NOT NULL,
			password VARCHAR(120) NOT NULL,
			created_at DATETIME,
			PRIMARY KEY (id)
		);
	`
	_, err := db.Exec(queryAuth)
	if err != nil {
		log.Fatalf("CreateDatabase err %s", err)
	}
	queryBooks := `
		CREATE TABLE books (
			id INT AUTO_INCREMENT,
			name VARCHAR(80) NOT NULL,
			author VARCHAR(80) NOT NULL,
			category VARCHAR(50),
			image VARCHAR(300),
			PRIMARY KEY (id)
		);
	`
	_, err2 := db.Exec(queryBooks)
	if err2 != nil {
		log.Fatalf("CreateDatabase err %s", err)
	}
}
