package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"

	"github.com/hukurou-s/GWMailer/db/models"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
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

	db, _ := gorm.Open("postgres", "user=LEO dbname=gwmailer-db password='' sslmode=disable")
	defer db.Close()

	user := models.User{
		Name:     name,
		Password: password, // to hash
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

	db, err := gorm.Open("postgres", "user=LEO dbname=gwmailer-db password='' sslmode=disable")
	defer db.Close()

	if err != nil {
		fmt.Print(err)
	}

	user := models.User{}
	db.First(&user, "name = ?", name)

	if user.Password != password {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return c.Render(http.StatusOK, "mypage", map[string]interface{}{
		"UserName": name,
	})
}
