package html

import (
	"net/http"

	"github.com/faagerholm/page/html/blog"
	"github.com/labstack/echo/v4"
)

func Blog(c echo.Context) error {
	content := struct {
		params
		Posts []blog.Post
	}{
		params: params{
			Title: "Blog",
		},
		Posts: blog.GetPosts(),
	}
	return c.Render(http.StatusOK, "blog", content)
}
