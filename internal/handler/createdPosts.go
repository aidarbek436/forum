package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/aidarbek436/forum/internal/repository"
)

func (h *handler) CreatedPostsHandle(w http.ResponseWriter, r *http.Request) { // DONE
	if r.URL.Path != "/createdposts" {
		err := ErrorPage(w, 404, "404 Page not Found: /createdposts")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	if r.Method != "POST" {
		err := ErrorPage(w, 405, "Method not Allowed:Should be Post")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	username, ok := GetUser(r.Context()).(string)
	if !ok {
		err := ErrorPage(w, 401, "Unauhtorized:session does not exist")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	created_posts, err := h.storage.GetCreatedPosts(username)
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error:GetCreatedPosts")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	var postsAndName repository.PostsWithName
	postsAndName.Post = created_posts
	postsAndName.Name = username
	absPath, err := filepath.Abs("front/root.html")
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error:filepath")
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
	err = temp.Execute(w, postsAndName)
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error:execute template")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}
