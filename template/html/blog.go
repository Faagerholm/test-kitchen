package html

import (
	"net/http"

	"github.com/faagerholm/page/auth"
	"github.com/faagerholm/page/html/blog"
	"github.com/faagerholm/page/session"
	"github.com/labstack/echo/v4"
)

func Blog(c echo.Context) error {
	content := struct {
		params
		Posts []blog.Post
	}{
		params: params{
			Title: "Blog",
			User:  auth.GetUser(session.ID(c.Request())),
		},
		Posts: blog.GetPosts(),
	}
	return c.Render(http.StatusOK, "blog", content)
}
