package entity

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
