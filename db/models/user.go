package models

import (
	"time"
)

type User struct {
	Id        int64
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
