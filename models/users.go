package models

import "time"

type User struct {
	ID         int
	Username   string
	Name       string
	ProfilePic string
	Email      string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
