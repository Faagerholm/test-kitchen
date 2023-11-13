package html

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/faagerholm/page/auth"
	"github.com/faagerholm/page/session"
	"github.com/labstack/echo/v4"
)

func MapPage(c echo.Context) error {
	p := params{
		Title: "Map",
		User:  auth.GetUser(session.ID(c.Request())),
	}

	return c.Render(http.StatusOK, "map", p)
}

var messageChannels = make(map[string]chan string)

func disconnect(id string) {
	close(messageChannels[id])
	delete(messageChannels, id)
	log.Debug(id + ": channel deleted")
}

func broadcastMessage(message string) {
	for id, ch := range messageChannels {
		log.Debug("Sending map location to channel",
			slog.String("coordinates", message),
			slog.String("ch", id))
		ch <- message
	}
}

func MapDrawer(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	id := session.ID(c.Request())

	ch := make(chan string)
	defer close(ch)

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()
	messageChannels[id] = ch

	template := `
message: %s
data: %s
`
	for {
		select {
		case <-ctx.Done():
			disconnect(id)
			return nil
		case msg := <-ch:
			log.Debug("Got message", slog.String("msg", msg), slog.String("ch", id))
			err := c.String(http.StatusOK, fmt.Sprintf(template, "map", msg))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			c.Response().Flush()
		}
	}
}

func generateRandomLatitude(r *rand.Rand) float64 {
	// Finland's latitude range is approximately 60.0 to 70.0
	return 60.0 + r.Float64()*(70.0-60.0)
}

func generateRandomLongitude(r *rand.Rand) float64 {
	// Finland's longitude range is approximately 20.0 to 31.0
	return 20.0 + r.Float64()*(31.0-20.0)
}

func StreamRandomLocation() {
	// Seed the random number generator with the current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// Generate random latitude and longitude in Finland
		latitude := generateRandomLatitude(r)
		longitude := generateRandomLongitude(r)

		message := fmt.Sprintf(`{"lat": %.4f, "lng": %.4f}`, latitude, longitude)
		broadcastMessage(message)
		time.Sleep(3 * time.Second)
	}
}
