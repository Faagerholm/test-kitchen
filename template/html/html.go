package html

import (
	"errors"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/faagerholm/page/auth"
	"github.com/faagerholm/page/session"
	"github.com/labstack/echo/v4"
)

type params struct {
	Title string
	User  *auth.User
}

var tpl *template.Template

func InitTemplates() {
	tpl = template.New("root")
	tpl = template.Must(tpl.ParseGlob("html/templates/*.html"))
}

type Template struct {
	templates *template.Template
}

func NewRenderer() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("html/templates/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Index(c echo.Context) error {
	d := struct {
		params
		Time string
	}{
		params: params{Title: "Welcome"},
		Time:   time.Now().Format(time.Kitchen),
	}
	return c.Render(http.StatusOK, "index", d)
}

func LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func Login(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to parse login form")
	}
	cookie, err := auth.Login(
		session.ID(c.Request()),
		auth.LoginForm{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, auth.UserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/")
}

func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "secret",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	return c.Redirect(http.StatusOK, "/login")
}

func KitchenTime(c echo.Context) error {
	t := time.Now().Format(time.Kitchen)
	return c.String(http.StatusOK, t)
}
