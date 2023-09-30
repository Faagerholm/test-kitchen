package html

import (
	"embed"
	"net/http"
	"text/template"
	"time"
)

//go:embed *
var files embed.FS

var (
	index    = parse("index.html")
	login    = parse("login.html")
	notFound = parse("404.html")
)

func Index(w http.ResponseWriter, r *http.Request) error {
	params := struct {
		Title string
		Time  string
	}{
		Title: "Helloo..",
		Time:  time.Now().Format(time.Kitchen),
	}
	return index.ExecuteTemplate(w, "layout.html", params)
}

func Login(w http.ResponseWriter, r *http.Request) error {
	return login.ExecuteTemplate(w, "layout.html", nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotFound)
	return notFound.ExecuteTemplate(w, "layout.html", nil)
}

func parse(page string) *template.Template {
	return template.Must(
		template.New("index").ParseFS(files, "layout.html", page),
	)
}
