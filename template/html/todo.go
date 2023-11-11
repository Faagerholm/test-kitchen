package html

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/faagerholm/page/logger"
	"github.com/labstack/echo/v4"
)

var log = slog.New(logger.NewHandler(&slog.HandlerOptions{
	Level:       slog.LevelDebug,
	AddSource:   false,
	ReplaceAttr: nil,
}))

type Task struct {
	ID   int
	Text string
}

type board struct {
	params
	New        []Task
	InProgress []Task
	Done       []Task
}

var Board = board{
	params: params{
		Title: "Todos",
	},
	New: []Task{
		{1, "Test out HTMX"},
		{2, "Demo HTMX"},
	},
	InProgress: []Task{
		{3, "Work on html/template"},
	},
	Done: []Task{
		{4, "Setup test-kitchen"},
	},
}

func TodoPage(c echo.Context) error {
	return c.Render(http.StatusOK, "board", Board)
}

func TodoAdd(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Could not parse input form")
	}
	todo := c.FormValue("new-todo")
	Board.New = append(Board.New, Task{len(Board.New) + 1, todo})
	d := struct {
		ID   int
		Text string
	}{len(Board.New), todo}
	return c.Render(http.StatusOK, "card", d)
}

func TodoMove(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Could not parse move input form")
	}
	task := c.FormValue("task")
	from := c.FormValue("from")
	to := c.FormValue("to")

	id, err := strconv.Atoi(task)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to convert task id to int")
	}

	var moving Task
	switch from {
	case "new":
		var list []Task
		moving, list = extractMovingTask(id, Board.New)
		Board.New = list
	case "in-progress":
		moving, Board.InProgress = extractMovingTask(id, Board.InProgress)
	case "done":
		moving, Board.Done = extractMovingTask(id, Board.Done)
	default:
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("%s not a valid move source", from),
		)
	}

	switch to {
	case "new":
		Board.New = append(Board.New, moving)
	case "in-progress":
		Board.InProgress = append(Board.InProgress, moving)
	case "done":
		Board.Done = append(Board.Done, moving)
	default:
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("%s not a valid move target", to),
		)
	}

	log.Debug(fmt.Sprintf("moving Task from %s to %s", from, to),
		slog.Any("Task", moving),
		slog.Any("board", Board),
	)
	return c.String(http.StatusOK, "Moved!")
}

func extractMovingTask(ID int, tasks []Task) (Task, []Task) {
	list := make([]Task, 0, len(tasks)-1)
	var task Task
	for _, t := range tasks {
		if ID == t.ID {
			task = t
		} else {
			list = append(list, t)
		}
	}
	return task, list
}
