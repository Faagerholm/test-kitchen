package auth

import (
	"fmt"
	"net/http"
)

var UserNotFound = fmt.Errorf("User not Found")

const SecretCookie = "chocolate-cookie"

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(f LoginForm) (*http.Cookie, error) {
	if f.Username == "Albert" && f.Password == "secret" {
		return &http.Cookie{
			Name:  "secret",
			Value: SecretCookie,
		}, nil
	}

	return nil, UserNotFound
}
