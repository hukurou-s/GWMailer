package models

import (
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model
	Address  string `gorm:"size:255;not null;unique"`
	Password string
	Server   string
	UserID   uint
}
