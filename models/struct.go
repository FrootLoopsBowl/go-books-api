package models

import "time"

type Books struct {
	Id       int
	Name     string
	Author   string
	Image    string
	Category string
}

type User struct {
	Username  string
	Password  string
	CreatedAt time.Time
}
