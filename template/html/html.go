package html

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/faagerholm/page/auth"
	"github.com/faagerholm/page/html/blog"
	"github.com/faagerholm/page/session"
)

type baseParams struct {
	PageTitle string
	User      *auth.User
}

var tpl *template.Template

func InitTemplates() {
	tpl = template.New("root")
	tpl = template.Must(tpl.ParseGlob("html/templates/*.html"))
}

func Index(w http.ResponseWriter, r *http.Request) error {
	p := struct {
		baseParams
		Time string
	}{
		Time: time.Now().Format(time.Kitchen),
	}
	p.User = auth.GetUser(session.ID(r))
	p.PageTitle = "Hello..."

	return tpl.ExecuteTemplate(w, "index.html", p)
}

func Blog(w http.ResponseWriter, r *http.Request) error {
	p := struct {
		baseParams
		Posts []blog.Post
	}{
		Posts: blog.GetPosts(),
	}
	p.User = auth.GetUser(session.ID(r))
	p.PageTitle = "Blog"
	return tpl.ExecuteTemplate(w, "blog.html", p)
}

func Login(w http.ResponseWriter, r *http.Request) error {
	var p baseParams

	return tpl.ExecuteTemplate(w, "login.html", p)
}

func NotFound(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	return tpl.ExecuteTemplate(w, "404.html", nil)
}

type Task struct {
	Title string
}

type board struct {
	New        []Task
	InProgress []Task
	Done       []Task
}

var Board = board{
	New: []Task{
		{"Test out HTMX"},
		{"Demo HTMX"},
	},
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := r.Form.Get("new-todo")
	Board.New = append(Board.New, Task{t})
	card := fmt.Sprintf(`<div id="done-%d" class="button is-fullwidth mt-1" draggable="true">
          %s
        </div>`, len(Board.New), t)
	fmt.Fprint(w, card)
}

func MoveTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var p struct {
		baseParams
		Board board
	}
	p.PageTitle = "Todo"

	to := r.Form.Get("to")
	task := r.Form.Get("task")

	if task == "" || to == "" {
		if err := tpl.ExecuteTemplate(w, "board.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	slog.Info("current board", "task", task, "to", to, "b", Board)
	s := strings.Split(task, "-")
	i, _ := strconv.Atoi(s[1])

	var old Task
	switch s[0] {
	case "new":
		old = Board.New[i]
		Board.New = append(Board.New[:i], Board.New[i+1:]...)
	case "progress":
		old = Board.InProgress[i]
		Board.InProgress = append(Board.InProgress[:i], Board.InProgress[i+1:]...)
	case "done":
		old = Board.Done[i]
		Board.Done = append(Board.Done[:i], Board.Done[i+1:]...)
	}
	switch to {
	case "todo-new":
		Board.New = append(Board.New, old)
	case "todo-progress":
		Board.InProgress = append(Board.InProgress, old)
	case "todo-done":
		Board.Done = append(Board.Done, old)
	}

	slog.Info("Update todos", "task", s, "index", i, "board", Board)
	p.Board = Board

	err = tpl.ExecuteTemplate(w, "board.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Todos(w http.ResponseWriter, r *http.Request) error {
	var p struct {
		baseParams
		Board board
	}
	p.PageTitle = "Todo"
	p.User = auth.GetUser(session.ID(r))
	p.Board = Board
	return tpl.ExecuteTemplate(w, "todo.html", p)
}

type CounterParams struct {
	Global, Session int
}

func GlobalCounter(w http.ResponseWriter, r *http.Request, p CounterParams) error {
	return tpl.ExecuteTemplate(w, "counter-global.html", p)
}

func SessionCounter(w http.ResponseWriter, r *http.Request, p CounterParams) error {
	return tpl.ExecuteTemplate(w, "counter-session.html", p)
}

func Counter(w http.ResponseWriter, r *http.Request, p CounterParams) error {
	return tpl.ExecuteTemplate(w, "counter.html", p)
}
