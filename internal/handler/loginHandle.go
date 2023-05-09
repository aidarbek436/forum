package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func (h *handler) LoginHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		err := ErrorPage(w, 404, "404 Page not Found:/login")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Wrong Url")
		return
	}
	switch r.Method {
	case http.MethodGet:
		_, err := r.Cookie("forum")
		if err == nil {
			err := ErrorPage(w, 400, "Bad request:session already exists")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		absPath, err := filepath.Abs("front/login.html")
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:FilePath")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		temp, err := template.ParseFiles(absPath)
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:Parse template")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		err = temp.Execute(w, nil)
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:tempalte execute")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
	case http.MethodPost:
		_, err := r.Cookie("forum")
		if err == nil {
			err := ErrorPage(w, 400, "Bad request:session already exists")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		username := r.FormValue("username")
		usernames := strings.ReplaceAll(username, " ", "")
		if usernames == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if HtmlInjectionCheck(username) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		password := r.FormValue("password")
		passwords := strings.ReplaceAll(password, " ", "")
		if passwords == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if HtmlInjectionCheck(password) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		userPassword, err := h.storage.GetUser(username)
		if err != nil {
			fmt.Println("err with login GetUser:", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if password != userPassword {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		setSession(username, password, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return

	default:
		err := ErrorPage(w, 405, "Method not Allowed:should be get or post")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	return
}
