package repository

import (
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}

type Post struct {
	Id         int
	Title      string
	Content    string
	Like       int
	Dislike    int
	Author     string
	IdOfAuthor int
	C_sport    int
	C_history  int
	C_politics int
	C_science  int
	C_art      int
}

type Comment struct {
	Id       int
	Author   string
	Text     string
	Like     int
	Dislike  int
	IdOfPost int
}

type PostsWithName struct {
	Post []Post
	Name string
}

type PostWithComments struct {
	Post     Post
	Comments []Comment
	Name     string
}

type Storage interface {
	UserIsExist(string) bool
	EmailIsExist(string) bool
	PostUser(User) error
	GetUser(string) (string, error)
	PostPost(Post) error
	GetAllPosts() ([]Post, error)
	GetCategoriesPosts(string) ([]Post, error)
	GetPost(int) (Post, error)
	GetCreatedPosts(string) ([]Post, error)
	GetComments(int) ([]Comment, error)
	PostComment(Comment) error

	UpdateLikePost(int, int, int) error
	UpdateDislikePost(int, int, int) error
	PostIsLike(string, int) (int, bool)
	InsertUserLike(int, string, int) error
	UpdateIsLikePost(string, int, int) error

	UpdateLikeComment(int, int, int) error
	UpdateDislikeComment(int, int, int) error
	CommentIsLike(string, int, int) (int, bool)
	InsertUserLikeComment(int, string, int, int) error
	UpdateIsLikeComment(string, int, int, int) error

	GetLikedPosts(string) ([]int, error)
}
