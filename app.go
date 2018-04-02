package main

import (
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
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
	}

	db_user = os.Getenv("USER_NAME")
	db_name = os.Getenv("DB_NAME")
	db_password = os.Getenv("DB_PASSWORD")
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
	e.GET("/login", getLogin)
	e.GET("/user/new", getUserNew)
	e.POST("/user/create", postCreateUser)
	e.POST("/mypage", postMypage)
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
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	db, _ := gorm.Open("postgres", "user="+db_user+" dbname="+db_name+" password='"+db_password+"' sslmode=disable")
	defer db.Close()

	user := models.User{
		Name:     name,
		Password: toHash(password),
	}

	db.Create(&user)

	return c.Redirect(http.StatusSeeOther, "/login")
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", map[string]interface{}{})
}

func postMypage(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")

	db, err := gorm.Open("postgres", "user="+db_user+" dbname="+db_name+" password='"+db_password+"' sslmode=disable")

	defer db.Close()

	if err != nil {
		fmt.Print(err)
	}

	user := models.User{}
	db.First(&user, "name = ?", name)

	if user.Password != toHash(password) {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	session := session.Default(c)
	session.Set("USERID", user.ID)
	session.Save()

	return c.Render(http.StatusOK, "mypage", map[string]interface{}{
		"UserName": name,
	})
}

func toHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
