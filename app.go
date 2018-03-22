package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
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
	e.POST("/mypage", postMypage)
	e.Logger.Fatal(e.Start(":1323"))
}

func getIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", map[string]interface{}{
		"PageName": "Top Page",
	})
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", map[string]interface{}{})
}

func postMypage(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")

	if name != "hoge" || password != "piyo" {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	return c.Render(http.StatusOK, "mypage", map[string]interface{}{
		"UserName": name,
	})
}
