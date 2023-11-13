package html

import (
	"context"
	"fmt"
	"net/http"

	"github.com/faagerholm/page/session"
	"github.com/faagerholm/page/store"
	"github.com/labstack/echo/v4"
)

func IncrementCounter(c echo.Context) error {
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch {
	case c.QueryParams().Has("global"):
		store.IncrementGlobal()
	case c.QueryParams().Has("session"):
		value, err := store.IncrementSession(sessionID.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, fmt.Sprintf("%d", value))
	}

	return nil
}

func CounterGet(c echo.Context) error {
	sessionID, err := c.Cookie("sessionID")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	global, session := store.Get(sessionID.Value)
	data := struct {
		Global  int
		Session int
	}{global, session}
	return c.Render(http.StatusOK, "counter", data)
}

var EventMsg = `
event: %s
data: %d
`

func CounterEvent(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	// Create a context to handle client disconnections
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	id := session.ID(c.Request())
	ch := session.ConnectClient(id)
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			log.Info("Client disconnected")
			session.DisconnectClient(id)
			return nil
		case m := <-ch:
			log.Info("got update", "message", m)
			err := c.String(http.StatusOK, fmt.Sprintf(EventMsg, "counter", m))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			c.Response().Flush()
		}
	}
}
