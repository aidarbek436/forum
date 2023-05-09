package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/aidarbek436/forum/internal/repository"
)

func (h *handler) ShowpostHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/showpost/" {
		err := ErrorPage(w, 404, "404 Page not Found:/showpost/")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Wrong Url")
		return
	}

	post_string_id := r.URL.Query().Get("id")
	post_string_ids := strings.ReplaceAll(post_string_id, " ", "")
	if post_string_ids == "" {
		err := ErrorPage(w, 400, "400 Bad Request")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	post_id, err := strconv.Atoi(post_string_id)
	if err != nil {
		err := ErrorPage(w, 400, "400 Bad Request")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	switch r.Method { // todo 1)handle url query   2) handle http errors
	case "GET":
		username, ok := GetUser(r.Context()).(string)
		if !ok {
			fmt.Println("session does not exist showpost get")
		}
		fmt.Println(post_id)
		post, err := h.storage.GetPost(post_id)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:GetPost")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		if post == (repository.Post{}) {
			err := ErrorPage(w, 400, "400 Bad Request")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		comments, err := h.storage.GetComments(post_id)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:GetComments")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
		fmt.Println(comments)
		var postWithComments repository.PostWithComments
		postWithComments.Post = post
		postWithComments.Comments = comments
		postWithComments.Name = username
		absPath, err := filepath.Abs("front/showpost.html")
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

		err = temp.Execute(w, postWithComments)
		if err != nil {
			fmt.Println(err)
			err := ErrorPage(w, 500, "Internal server error:Execute Template")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Internal server error")
			return
		}
	case "POST":
		username, ok := GetUser(r.Context()).(string)
		if !ok {
			fmt.Println("session does not exist created posts method post")
			fmt.Println(ok)
			err := ErrorPage(w, 401, "Unauhtorized:session does not exist")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		fmt.Println(username)
		num, flag := h.storage.PostIsLike(username, post_id)
		if !flag {
			err = h.storage.InsertUserLike(0, username, post_id)
			if err != nil {
				err := ErrorPage(w, 500, "Internal server error:PostIsLike")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Internal server error")
				return
			}
		}
		like := r.FormValue("postLike")
		switch like {
		case "like":
			if num == 0 {
				err = h.storage.UpdateLikePost(post_id, 1, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikePost(username, post_id, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == 1 {
				err = h.storage.UpdateLikePost(post_id, -1, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikePost(username, post_id, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == -1 {
				err = h.storage.UpdateLikePost(post_id, 1, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikePost(username, post_id, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			}
		case "dislike":
			if num == 0 {
				err = h.storage.UpdateDislikePost(post_id, 0, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikePost(username, post_id, -1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == 1 {
				err = h.storage.UpdateDislikePost(post_id, 1, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikePost(username, post_id, -1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == -1 {
				err = h.storage.UpdateDislikePost(post_id, 0, -1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikePost(username, post_id, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikePost")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			}
		case "commentLike":
			comment_id_string := r.FormValue("commentId")
			comment_id_strings := strings.ReplaceAll(comment_id_string, " ", "")
			if comment_id_strings == "" {
				err := ErrorPage(w, 400, "400 Bad Request:comment_id")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			comment_id, err := strconv.Atoi(comment_id_string)
			if err != nil {
				err := ErrorPage(w, 400, "400 Bad Request:comment_id")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			num, flag := h.storage.CommentIsLike(username, post_id, comment_id)
			if !flag {
				err = h.storage.InsertUserLikeComment(0, username, post_id, comment_id)
				if err != nil {
					fmt.Println("InsertUserLikeComment")
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			}
			if num == 0 {
				err = h.storage.UpdateLikeComment(comment_id, 1, 0)
				if err != nil {
					fmt.Println("UpdateLikeComment")
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikeComment(username, post_id, comment_id, 1)
				if err != nil {
					fmt.Println("UpdateIsLikeComment")
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == 1 {
				err = h.storage.UpdateLikeComment(comment_id, -1, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikeComment(username, post_id, comment_id, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == -1 {
				err = h.storage.UpdateLikeComment(comment_id, 1, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikeComment(username, post_id, comment_id, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:LikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			}

		case "commentDislike":
			comment_id_string := r.FormValue("commentId")
			comment_id_strings := strings.ReplaceAll(comment_id_string, " ", "")
			if comment_id_strings == "" {
				err := ErrorPage(w, 400, "400 Bad Request:comment_id")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			comment_id, err := strconv.Atoi(comment_id_string)
			if err != nil {
				err := ErrorPage(w, 400, "400 Bad Request:comment_id")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			num, flag := h.storage.CommentIsLike(username, post_id, comment_id)
			if !flag {
				err = h.storage.InsertUserLikeComment(0, username, post_id, comment_id)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			}
			if num == 0 {
				err = h.storage.UpdateDislikeComment(comment_id, 0, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikeComment(username, post_id, comment_id, -1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == 1 {
				err = h.storage.UpdateDislikeComment(comment_id, 1, 1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikeComment(username, post_id, comment_id, -1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			} else if num == -1 {
				err = h.storage.UpdateDislikeComment(comment_id, 0, -1)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
				err = h.storage.UpdateIsLikeComment(username, post_id, comment_id, 0)
				if err != nil {
					err := ErrorPage(w, 500, "Internal server error:DislikeComment")
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Internal server error")
					return
				}
			}
		default:
			var commentInput repository.Comment
			comment := r.FormValue("content")
			comments := strings.ReplaceAll(comment, " ", "")
			if comments == "" {
				err := ErrorPage(w, 400, "400 Bad Request:nil comment value")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			if HtmlInjectionCheck(comment) == false {
				err := ErrorPage(w, 400, "400 Bad Request: html injection")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			commentInput.Author = username
			commentInput.Text = comment
			commentInput.IdOfPost = post_id
			if err := h.storage.PostComment(commentInput); err != nil {
				fmt.Println("err with PostComment", err)
				fmt.Println(err)
				err := ErrorPage(w, 500, "Internal server error:PostComment")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Internal server error")
				return
			}
		}
		http.Redirect(w, r, "/showpost/?id="+post_string_id, http.StatusFound)

	default:
		fmt.Println("wrong method")
		err := ErrorPage(w, 405, "Method not Allowed:should be Get or Post")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	return
}
