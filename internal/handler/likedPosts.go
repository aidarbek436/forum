package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/aidarbek436/forum/internal/repository"
)

func (h *handler) LikedPostsHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likedposts" {
		err := ErrorPage(w, 404, "404 Page not Found:/likedposts")
		if err != nil {
			fmt.Println(err)
		}
		return
	} // DONE
	if r.Method != "POST" {
		err := ErrorPage(w, 405, "Method not Allowed:should be Post")
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
	Liked_posts, err := h.storage.GetLikedPosts(username)
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error:GetLikedPosts")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	var created_posts []repository.Post
	for i := 0; i < len(Liked_posts); i++ {
		post := repository.Post{}
		post, err = h.storage.GetPost(Liked_posts[i])
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:GetPost")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		created_posts = append(created_posts, post)

	}
	var postsAndName repository.PostsWithName
	postsAndName.Post = created_posts
	postsAndName.Name = username
	absPath, err := filepath.Abs("front/root.html")
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error:FilePath")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	temp, err := template.ParseFiles(absPath)
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error: template parse")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	err = temp.Execute(w, postsAndName)
	if err != nil {
		err := ErrorPage(w, 500, "Internal server error:template execute")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}
