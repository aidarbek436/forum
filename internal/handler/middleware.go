package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func middleware(next http.HandlerFunc) http.HandlerFunc { // TODO: Multiple sessions handle
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := SetUser(r.Context(), nil)
		c, err := r.Cookie("forum")
		if err == nil {
			sessionToken := c.Value

			userSession, exists := sessions[sessionToken]
			if !exists {
				http.SetCookie(w, &http.Cookie{
					Name:    "forum",
					Value:   "",
					Expires: time.Now(),
					Path:    "/",
				})
				next(w, r)
				return
			}
			if userSession.isExPired() {
				fmt.Println("session is expired")
				delete(sessions, sessionToken)
				http.SetCookie(w, &http.Cookie{
					Name:    "forum",
					Value:   "",
					Expires: time.Now(),
				})
				err := ErrorPage(w, 401, "Unauhtorized:session is expired")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			if len(sessions) > 1 {
				for key, elements := range sessions {
					if userSession.username == elements.username {
						delete(sessions, key)
					}
				}
			}
			expireAt := time.Now().Add(900 * time.Second)

			sessions[sessionToken] = session{
				username: userSession.username,
				expiry:   expireAt,
			}
			http.SetCookie(w, &http.Cookie{
				Name:    "forum",
				Value:   sessionToken,
				Expires: expireAt,
				Path:    "/",
			})
			ctx = SetUser(ctx, userSession.username)

			next(w, r.WithContext(ctx))
			return
		} else {
			switch r.Method {
			case "POST":
				if r.URL.Path == "/register" {
					next(w, r.WithContext(ctx))
					return
				}
				if r.URL.Path == "/login" {
					next(w, r.WithContext(ctx))
					return
				}
				if r.URL.Path == "/" {
					next(w, r.WithContext(ctx))
					return
				}
				err := ErrorPage(w, 401, "Unauhtorized:session does not exist")
				if err != nil {
					fmt.Println(err)
				}
				return
			default:
				next(w, r.WithContext(ctx))
				return
			}
		}
	})
}

func SetUser(ctx context.Context, value interface{}) context.Context {
	return context.WithValue(ctx, "user", value)
}
