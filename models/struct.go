package models

import "time"

type Books struct {
	name     string
	author   string
	image    string
	category string
}

type User struct {
	Username  string
	Password  string
	CreatedAt time.Time
}
