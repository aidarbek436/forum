package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

type Credentials struct {
	username string
	password string
}

func (s session) isExPired() bool {
	return s.expiry.Before(time.Now())
}

func setSession(username string, password string, w http.ResponseWriter) {
	cred := Credentials{username, password}
	sessionToken := uuid.NewString()
	expireAt := time.Now().Add(900 * time.Second)
	sessions[sessionToken] = session{
		username: cred.username,
		expiry:   expireAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "forum",
		Value:   sessionToken,
		Expires: expireAt,
	})
}

func deleteSession(sessionToken string, w http.ResponseWriter) {
	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "forum",
		Value:   sessionToken,
		Expires: time.Now(),
	})
}
