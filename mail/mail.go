package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/mvader/go-imapreader.v1"
)

var (
	imap_server string
	user_name   string
	password    string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	imap_server = os.Getenv("IMAP_SERVER")
	user_name = os.Getenv("MAIL_ADDRESS")
	password = os.Getenv("MAIL_PASSWORD")
}

func main() {

	r, err := imapreader.NewReader(imapreader.Options{
		Addr:     imap_server,
		Username: user_name,
		Password: password,
		TLS:      true,
		Timeout:  60 * time.Second,
		MarkSeen: true,
	})

	if err != nil {
		panic(err)
	}

	if err := r.Login(); err != nil {
		panic(err)
	}
	defer r.Logout()
}
