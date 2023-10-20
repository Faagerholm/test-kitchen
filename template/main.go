package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/faagerholm/page/auth"
	"github.com/faagerholm/page/html"
	"github.com/faagerholm/page/observer"
	"github.com/faagerholm/page/session"
	"github.com/faagerholm/page/store"
)

func main() {
	s := http.NewServeMux()

	s.HandleFunc("/", index)
	s.HandleFunc("/login", login)
	s.HandleFunc("/logout", logout)
	s.HandleFunc("/blog", blog)

	s.HandleFunc("/todo", todo)
	s.HandleFunc("/todo/new", html.AddTodo)
	s.HandleFunc("/todo/move", html.MoveTodo)

	s.HandleFunc("/counter", counter)
	s.HandleFunc("/counter/event", counterEvent)
	s.HandleFunc("/time", clock)

	ob := observer.NewMiddleware(s)
	sh := session.NewMiddleware(ob.Next, session.WithSecure(false))

	session.NewBroadcaster()
	html.InitTemplates()
	store.NewCounter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: sh,
		//		ReadTimeout:  time.Second * 10,
		//		WriteTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on %v\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if err := html.NotFound(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	if err := html.Index(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func clock(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.Kitchen)
	fmt.Fprint(w, t)
}

func blog(w http.ResponseWriter, r *http.Request) {
	if err := html.Blog(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func todo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		html.AddTodo(w, r)
		return
	}
	if err := html.Todos(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func counterEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a context to handle client disconnections
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := session.ID(r)
	ch := session.ConnectClient(id)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Client disconnected")
			session.DisconnectClient(session.ID(r))
			return
		case m := <-ch:
			slog.Info("got update", "message", m)
			fmt.Fprintf(w, "event: globalCounter\ndata: %d\n\n", m)
			w.(http.Flusher).Flush()
		}
	}
}

// counter accepts Post, to increment counter
func counter(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("sessionID")
	if err != nil {
		// Session should be set by our session-middleware
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	var global, session int
	switch r.Method {
	case http.MethodGet:
		// yes
	case http.MethodPost:
		v := r.URL.Query()
		switch {
		case v.Has("global"):
			store.IncrementGlobal()
			return
		case v.Has("session"):
			count, err := store.IncrementSession(c.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "%d", count)
			return
		default:
			http.Error(w, "invalid type for counter", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "unsupported method", http.StatusBadRequest)
		return
	}
	global, session = store.Get(c.Value)
	err = html.Counter(w, r, html.CounterParams{Global: global, Session: session})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie, err := auth.Login(
			session.ID(r),
			auth.LoginForm{
				Username: r.Form.Get("username"),
				Password: r.Form.Get("password"),
			})
		if err != nil {
			switch {
			case errors.Is(err, auth.UserNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusFound)

		return
	}
	if err := html.Login(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "secret",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("secret")
		if err != nil {
			slog.Error("cookie", slog.Any("error", err))
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if cookie == nil || cookie.Value != auth.SecretCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
