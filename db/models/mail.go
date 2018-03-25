package models

import (
	"time"
)

type Mail struct {
	Id        int64
	From      string
	To        string
	Cc        string
	Subject   string
	body      text
	GotAt     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
