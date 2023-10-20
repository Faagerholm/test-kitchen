package session

import (
	"net/http"

	"github.com/segmentio/ksuid"
)

type MiddlewareOpts func(*Middleware)

func NewMiddleware(next http.Handler, opts ...MiddlewareOpts) Middleware {
	mw := Middleware{
		Next:     next,
		Secure:   true,
		HTTPOnly: true,
	}

	for _, opt := range opts {
		opt(&mw)
	}

	return mw
}

type Middleware struct {
	Next     http.Handler
	Secure   bool
	HTTPOnly bool
}

func WithSecure(secure bool) MiddlewareOpts {
	return func(m *Middleware) {
		m.Secure = secure
	}
}

func WithHTTPOnly(httpOnly bool) MiddlewareOpts {
	return func(m *Middleware) {
		m.HTTPOnly = httpOnly
	}
}

func ID(r *http.Request) string {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (mw Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := ID(r)
	if id == "" {
		id = ksuid.New().String()
		http.SetCookie(w, &http.Cookie{Name: "sessionID", Value: id, Secure: mw.Secure, HttpOnly: mw.HTTPOnly})
	}

	mw.Next.ServeHTTP(w, r)
}
