package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aidarbek436/forum/internal/repository"
)

func (h *handler) CreatepostHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createpost" {
		err := ErrorPage(w, 404, "404 Page not Found:/createpost")
		if err != nil {
			fmt.Println(err)
		}
		return
	} // DONE
	switch r.Method {
	case http.MethodGet:
		_, ok := GetUser(r.Context()).(string)
		if !ok {
			err := ErrorPage(w, 401, "Unauhtorized:session does not exist")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		absPath, err := filepath.Abs("front/createpost.html")
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:filePath")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		temp, err := template.ParseFiles(absPath)
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:parse template")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err = temp.Execute(w, nil)
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:execute template")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	case http.MethodPost:
		username, ok := GetUser(r.Context()).(string)
		if !ok {
			err := ErrorPage(w, 401, "Unauhtorized:session does not exist")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		fmt.Println(username)
		var postINput repository.Post
		postINput.Author = username
		postINput.Title = r.FormValue("title")
		Titles := strings.ReplaceAll(postINput.Title, " ", "")
		if Titles == "" {
			err := ErrorPage(w, 400, "400 Bad Request: nil Title value")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if HtmlInjectionCheck(postINput.Title) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		postINput.Content = r.FormValue("content")
		Contents := strings.ReplaceAll(postINput.Content, " ", "")
		if Contents == "" {
			err := ErrorPage(w, 400, "400 Bad Request: nil Content value")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if HtmlInjectionCheck(postINput.Content) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		r.ParseForm()
		options := r.Form["option"]
		fmt.Println(options)
		if len(options) == 0 {
			err := ErrorPage(w, 400, "400 Bad Request: Absense of categories post should have AT LEAST ONE (x>0)")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		for _, option := range options {
			if HtmlInjectionCheck(option) == false {
				err := ErrorPage(w, 400, "400 Bad Request: html injection")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			if option == "art" {
				postINput.C_art = 1
			}
			if option == "history" {
				postINput.C_history = 1
			}
			if option == "politics" {
				postINput.C_politics = 1
			}
			if option == "science" {
				postINput.C_science = 1
			}
			if option == "sport" {
				postINput.C_sport = 1
			}
		}

		postINput.Like = 0
		postINput.Dislike = 0

		if err := h.storage.PostPost(postINput); err != nil {
			err := ErrorPage(w, 500, "Internal server error:PostPost")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		err := ErrorPage(w, 405, "Method not Allowed:should be get or post")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}
