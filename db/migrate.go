package main

import (
	"fmt"

	"github.com/hukurou-s/GWMailer/db/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "user=LEO dbname=gwmailer-db password='' sslmode=disable")
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("success")
	}

	db.AutoMigrate(&models.User{})
	//db.CreateTable(&models.User{})

}
