package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Mail struct {
	gorm.Model
	From       string `gorm:"size:255"`
	To         string `gorm:"size:255"`
	Cc         string `gorm:"size:255"`
	Subject    string `gorm:"size:255"`
	Body       string
	ReceivedAt time.Time
}
