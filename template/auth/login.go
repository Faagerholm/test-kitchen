package auth

import (
	"fmt"
	"net/http"
	"time"
)

var UserNotFound = fmt.Errorf("User not Found")

const SecretCookie = "chocolate-cookie"

type User struct {
	sessionID string
	Name      string
}

var Users []User

func GetUser(sessionID string) *User {
	for _, u := range Users {
		if sessionID == u.sessionID {
			return &u
		}
	}
	return nil
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(id string, l LoginForm) (*http.Cookie, error) {
	u := GetUser(id)
	if u != nil && u.Name != l.Username {
		return nil, fmt.Errorf("You already have an user with another name, please use that")
	}
	if u == nil {
		Users = append(Users, User{sessionID: id, Name: l.Username})
	}
	return &http.Cookie{
		Name:  "secret",
		Value: SecretCookie,
		// 24 hours
		Expires: time.Now().AddDate(0, 0, 1),
	}, nil
}
