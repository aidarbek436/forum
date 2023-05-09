package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/aidarbek436/forum/internal/repository"
)

func (h *handler) RootHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err := ErrorPage(w, 404, "404 Page not Found:/")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Wrong Url")
		return
	}
	switch r.Method {
	case http.MethodGet: // hande http errors
		username, ok := GetUser(r.Context()).(string)
		if !ok {
			fmt.Println("session does not exist root get")
		}

		fmt.Println(r.URL.Path)

		getAllPosts, err := h.storage.GetAllPosts()
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:GetAllPosts")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		var postsAndName repository.PostsWithName
		postsAndName.Post = getAllPosts
		postsAndName.Name = username
		absPath, err := filepath.Abs("front/root.html")
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:FilePath")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		temp, err := template.ParseFiles(absPath)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:Parse Temlate")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		err = temp.Execute(w, postsAndName)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:Execute Template")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
	case http.MethodPost:
		username, ok := GetUser(r.Context()).(string) // todo 1 handle http errors 2) handle r. body
		if !ok {
			fmt.Println(ok)
			fmt.Println("session does not exist")
		}
		fmt.Println(username)
		r.ParseForm()
		options := r.Form["option"]
		fmt.Println(options)
		var postsRes []repository.Post
		for _, option := range options {
			if HtmlInjectionCheck(option) == false {
				err := ErrorPage(w, 400, "400 Bad Request: html injection")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			posts, err := h.storage.GetCategoriesPosts(option)
			if err != nil {
				fmt.Println("err with roothandle getCategoryPosts", err)
				fmt.Println(err)
				err := ErrorPage(w, 500, "Internal server error:GetCategoriesPosts")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Internal server error")
				return
			}

			for _, p := range posts {
				fmt.Println("133")
				if len(postsRes) == 0 {
					fmt.Println("123")
					postsRes = append(postsRes, p)
				}
				for i, po := range postsRes {
					if p == po {
						break
					}
					if i == len(postsRes)-1 {
						postsRes = append(postsRes, p)
					}
				}
			}
			fmt.Println(postsRes, "1")
		}
		var postsAndName repository.PostsWithName
		postsAndName.Post = postsRes
		postsAndName.Name = username
		absPath, err := filepath.Abs("front/root.html")
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:FilePath")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		temp, err := template.ParseFiles(absPath)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:Parse Template")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		err = temp.Execute(w, postsAndName)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:Execute Template")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
	default:
		fmt.Println("wrong method")
		err := ErrorPage(w, 405, "Method not Allowed:should be Get or Post")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}

func GetUser(ctx context.Context) interface{} {
	return ctx.Value("user")
}
