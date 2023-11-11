package main

import (
	"net/http"

	"github.com/faagerholm/page/html"
	"github.com/faagerholm/page/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.New(false, false))

	e.Renderer = html.NewRenderer()

	e.GET("/", html.Index)
	e.GET("/counter", html.CounterGet)
	e.POST("/counter", html.IncrementCounter)
	e.GET("/events/counter", html.CounterEvent)
	e.GET("/time", html.KitchenTime)

	e.GET("/todo", html.TodoPage)
	e.POST("/todo/new", html.TodoAdd)
	e.PUT("/todo/move", html.TodoMove)

	e.GET("/blog", html.Blog)
	e.GET("/login", html.LoginPage)
	e.POST("/login", html.Login)
	e.GET("/logout", html.Logout)
	e.Static("/static", "assets")

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404", struct{ Title string }{"404"})
	}

	e.Logger.Fatal(e.Start(":8080"))
}
