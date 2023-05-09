package handler

import (
	"fmt"
	"net/http"
	"time"
)

func (h *handler) LogoutHandle(w http.ResponseWriter, r *http.Request) { // DONE
	if r.Method != "GET" {
		err := ErrorPage(w, 405, "Method not Allowed:should be Get")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	c, err := r.Cookie("forum")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	delete(sessions, sessionToken)
	fmt.Println(sessions)
	http.SetCookie(w, &http.Cookie{
		Name:    "forum",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
