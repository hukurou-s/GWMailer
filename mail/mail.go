package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	decode "github.com/curious-eyes/jmail"
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

	mails, err := r.List(imapreader.GMailInbox, imapreader.SearchUnseen)
	if err != nil {
		panic(err)
	}

	for _, mail := range mails {

		date := "Date: " + mail.Header["Date"][0] + "\n"
		from := "From: " + mail.Header["From"][0] + "\n"
		to := "To: " + mail.Header["To"][0] + "\n"
		subject := "Subject: " + mail.Header["Subject"][0] + "\n"
		contentType := "Content-Type: " + mail.Header["Content-Type"][0] + "\n"
		contentTransferEncoding := "Content-Transfer-Encoding: " + mail.Header["Content-Transfer-Encoding"][0] + "\n"
		message := date + from + to + subject + contentType + contentTransferEncoding + "\n" + string(mail.Body)

		r := strings.NewReader(message)
		m, _ := decode.ReadMessage(r)
		body, _ := m.DecBody()

		fmt.Printf("%s", body)
	}

}
