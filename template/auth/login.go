package auth

import (
	"fmt"
)

var UserNotFound = fmt.Errorf("User not Found")

const SecretCookie = "chocolate-cookie"

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(f LoginForm) (string, error) {
	if f.Username == "Albert" && f.Password == "secret" {
		return SecretCookie, nil
	}

	return "", UserNotFound
}
