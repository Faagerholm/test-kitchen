package main

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/faagerholm/page/auth"
	"github.com/faagerholm/page/html"
)

func main() {
	http.Handle("/", authMiddleware(loggingMiddleware(
		http.HandlerFunc(index),
	)))
	http.Handle("/login", loggingMiddleware(
		http.HandlerFunc(login),
	))
	http.Handle("/logout", loggingMiddleware(
		http.HandlerFunc(logout),
	))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		o := &observer{ResponseWriter: w}
		next.ServeHTTP(o, r)

		slog.Info("",
			slog.Int("status", o.status),
			slog.Duration("duration", time.Since(start)),
			slog.String("method", r.Method),
			slog.Any("endpoint", r.URL),
		)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie, err := auth.Login(auth.LoginForm{
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

		http.SetCookie(w, &http.Cookie{
			Name:  "secret",
			Value: cookie,
		})
		http.Redirect(w, r, "/", http.StatusFound)

		return
	}
	if cookie, _ := r.Cookie("secret"); cookie.String() == auth.SecretCookie {
		// user already authenticated
		http.Redirect(w, r, "/", http.StatusFound)
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

type observer struct {
	http.ResponseWriter
	status      int
	written     uint64
	wroteHeader bool
}

func (o *observer) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += uint64(n)

	return
}

func (o *observer) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
