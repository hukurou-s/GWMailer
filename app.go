package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/ipfans/echo-session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/hukurou-s/GWMailer/db/models"
	mailGetter "github.com/hukurou-s/GWMailer/mail"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var (
	db_user     string
	db_name     string
	db_password string
	secret_key  []byte
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
	}

	db_user = os.Getenv("USER_NAME")
	db_name = os.Getenv("DB_NAME")
	db_password = os.Getenv("DB_PASSWORD")
	key := os.Getenv("KEY")
	secret_key = convertTo32Byte(key)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	store := session.NewCookieStore([]byte("secret"))
	store.MaxAge(1200)
	e.Use(session.Sessions("GWMSESSION", store))

	e.GET("/", getIndex)
	e.GET("/logins/new", getNewLogin)
	e.POST("/logins", postLogins)
	e.GET("/user/new", getUserNew)
	e.POST("/user/create", postCreateUser)
	e.GET("/mypage", getMypage)
	e.GET("/accounts/new", getMailAccountNew)
	e.POST("/accounts/create", postCreateMailAccount)

	e.Logger.Fatal(e.Start(":1323"))
}

func getIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", map[string]interface{}{
		"PageName": "Top Page",
	})
}

func getUserNew(c echo.Context) error {
	return c.Render(http.StatusOK, "user_new", map[string]interface{}{})
}

func postCreateUser(c echo.Context) error {

	// registration db
	name := c.FormValue("name")
	password := c.FormValue("password")
	if password != c.FormValue("password_confirm") {
		return c.Redirect(http.StatusSeeOther, "/users/new")
	}

	db, _ := gorm.Open("postgres", "user="+db_user+" dbname="+db_name+" password='"+db_password+"' sslmode=disable")
	defer db.Close()

	user := models.User{
		Name:     name,
		Password: toHash(password),
	}

	db.Create(&user)

	return c.Redirect(http.StatusSeeOther, "/logins")
}

func getNewLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", map[string]interface{}{})
}

func postLogins(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")

	db, err := gorm.Open("postgres", "user="+db_user+" dbname="+db_name+" password='"+db_password+"' sslmode=disable")

	defer db.Close()

	if err != nil {
		fmt.Print(err)
	}

	user := models.User{}
	db.First(&user, "name = ?", name)

	if db.Find(&user, "name = ?", name).RecordNotFound() {
		return c.Redirect(http.StatusSeeOther, "/logins/new")
	}

	if user.Password != toHash(password) {
		return c.Redirect(http.StatusSeeOther, "/logins/new")
	}

	session := session.Default(c)
	session.Set("USERID", user.ID)
	session.Save()

	return c.Redirect(http.StatusSeeOther, "/mypage")

}

func getMypage(c echo.Context) error {

	session := session.Default(c)
	id := session.Get("USERID").(uint)

	db, err := gorm.Open("postgres", "user="+db_user+" dbname="+db_name+" password='"+db_password+"' sslmode=disable")

	defer db.Close()

	if err != nil {
		panic(err)
	}

	user := models.User{}
	db.First(&user, id)

	address := models.Address{}

	db.First(&address, "user_id = ?", id)

	address.Password = toDecrypt(address.Password)

	mailGetter.RegistUnSeenMail(address)

	mail := models.Mail{}
	db.Where("\"to\" = ?", address.Address).First(&mail)
	//db.First(&mail)

	return c.Render(http.StatusOK, "mypage", map[string]interface{}{
		"UserName":    user.Name,
		"MailAddress": address.Address,
		"Mail":        string(mail.From),
	})
}

func getMailAccountNew(c echo.Context) error {
	return c.Render(http.StatusOK, "mail_account_new", map[string]interface{}{})
}

func postCreateMailAccount(c echo.Context) error {
	// registration db
	session := session.Default(c)
	id := session.Get("USERID").(uint)

	mailAddress := c.FormValue("mail_address")
	password := c.FormValue("password")
	server := c.FormValue("server")

	if password != c.FormValue("password_confirm") {
		return c.Redirect(http.StatusSeeOther, "/accounts/new")
	}

	db, _ := gorm.Open("postgres", "user="+db_user+" dbname="+db_name+" password='"+db_password+"' sslmode=disable")
	defer db.Close()

	address := models.Address{
		Address:  mailAddress,
		Password: toEncrypt(password),
		Server:   server,
		UserID:   id,
	}

	db.Create(&address)

	return c.Redirect(http.StatusSeeOther, "/mypage")
}

func toHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func convertTo32Byte(key string) []byte {
	tmpKey := []byte(key)
	for len(tmpKey) < 32 {
		tmpKey = append(tmpKey, []byte(key)...)
	}
	secretKey := tmpKey[:32]
	return secretKey
}

func toEncrypt(password string) string {

	plainPassword := []byte(password)
	block, _ := aes.NewCipher(secret_key)

	cipherPassword := make([]byte, aes.BlockSize+len(plainPassword))
	iv := cipherPassword[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Printf("err: %s\n", err)
	}

	encryptStream := cipher.NewCTR(block, iv)

	encryptStream.XORKeyStream(cipherPassword[aes.BlockSize:], plainPassword)

	return hex.EncodeToString(cipherPassword)
}

func toDecrypt(password string) string {

	//plainPassword := []byte(password)
	cipherPassword, _ := hex.DecodeString(password)
	block, _ := aes.NewCipher(secret_key)

	decryptedPassword := make([]byte, len(cipherPassword[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherPassword[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedPassword, cipherPassword[aes.BlockSize:])

	return string(decryptedPassword)
}
