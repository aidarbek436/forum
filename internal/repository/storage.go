package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storage {
	return &storage{db: db}
}

func (s storage) PostUser(userInput User) error {
	if _, err := s.db.Exec("INSERT INTO users (Username, Email, Password) VALUES ($1,$2,$3)", userInput.Username, userInput.Email, userInput.Password); err != nil {
		return err
	}

	fmt.Println(userInput)

	return nil
}

func (s storage) GetUser(username string) (string, error) {
	stmt, err := s.db.Query("SELECT Password FROM users WHERE Username like '%" + username + "%'")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var password string

	for stmt.Next() {
		err = stmt.Scan(&password)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = stmt.Err()
	if err != nil {
		log.Fatal(err)
	}

	return password, nil
}

func (s storage) UserIsExist(username string) bool {
	err := s.db.QueryRow("SELECT Username FROM users WHERE Username=?", username).Scan(&username)
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func (s storage) EmailIsExist(email string) bool {
	err := s.db.QueryRow("SELECT Email FROM users WHERE Email=?", email).Scan(&email)
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func (s storage) PostPost(postInput Post) error {
	if _, err := s.db.Exec("INSERT INTO posts (Title, Content, Likes, Dislikes, Author, IdOfAuthor, C_sport, C_history, C_politics, C_science, C_art) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", postInput.Title,
		postInput.Content, postInput.Like, postInput.Dislike, postInput.Author, postInput.IdOfAuthor, postInput.C_sport, postInput.C_history, postInput.C_politics,
		postInput.C_science, postInput.C_art); err != nil {
		return err
	}

	return nil
}

func (s storage) GetAllPosts() ([]Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts")
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	Posts := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Author,
			&post.IdOfAuthor, &post.C_sport, &post.C_history, &post.C_politics, &post.C_science, &post.C_art)
		if err != nil {
			return nil, err
		}
		Posts = append(Posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

func (s storage) GetCategoriesPosts(category string) ([]Post, error) {
	stmt, err := s.db.Query("SELECT * FROM posts WHERE " + "C_" + category + " = 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	Posts := make([]Post, 0)

	for stmt.Next() {
		post := Post{}
		err = stmt.Scan(&post.Id, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Author, &post.IdOfAuthor, &post.C_art, &post.C_history, &post.C_politics, &post.C_science, &post.C_sport)
		if err != nil {
			return nil, err
		}
		Posts = append(Posts, post)
	}

	err = stmt.Err()
	if err != nil {
		log.Fatal(err)
	}
	return Posts, nil
}

func (s storage) GetPost(id int) (Post, error) {
	post := Post{}
	id_string := strconv.Itoa(id)
	stmt, err := s.db.Query("SELECT * FROM posts Where Id like '%" + id_string + "%'")
	if err != nil {
		return post, err
	}
	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&post.Id, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Author, &post.IdOfAuthor, &post.C_art, &post.C_history, &post.C_politics, &post.C_science, &post.C_sport)
		if err != nil {
			return post, err
		}
	}

	err = stmt.Err()
	if err != nil {
		log.Fatal(err)
	}
	return post, nil
}

func (s storage) GetCreatedPosts(username string) ([]Post, error) {
	rows, err := s.db.Query("SELECT * FROM posts WHERE Author like '%" + username + "%'")
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	Posts := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Like, &post.Dislike, &post.Author,
			&post.IdOfAuthor, &post.C_sport, &post.C_history, &post.C_politics, &post.C_science, &post.C_art)
		if err != nil {
			return nil, err
		}
		Posts = append(Posts, post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

func (s storage) GetComments(IdOfPost int) ([]Comment, error) {
	string_id_ofpost := strconv.Itoa(IdOfPost)
	rows, err := s.db.Query("SELECT * FROM comments WHERE IdOfPost like '%" + string_id_ofpost + "%'")
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	Comments := make([]Comment, 0)

	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(&comment.Id, &comment.Text, &comment.Like, &comment.Dislike, &comment.Author, &comment.IdOfPost)
		if err != nil {
			return nil, err
		}
		Comments = append(Comments, comment)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return Comments, nil
}

func (s storage) PostComment(CommentInput Comment) error {
	if _, err := s.db.Exec("INSERT INTO comments (Content, Likes, Dislikes, Author, IdOfPost) VALUES ($1,$2,$3,$4,$5)", CommentInput.Text,
		CommentInput.Like, CommentInput.Dislike, CommentInput.Author, CommentInput.IdOfPost); err != nil {
		return err
	}

	return nil
}

func (s storage) UpdateIsLikePost(username string, post_id, like int) error {
	if _, err := s.db.Exec("UPDATE posts_is_like SET IsLike = ? WHERE Username = ? AND Post_id = ?", like, username, post_id); err != nil {
		return err
	}
	return nil
}

func (s storage) UpdateLikePost(id_of_post, like, dislike int) error {
	string_id_ofpost := strconv.Itoa(id_of_post)
	if _, err := s.db.Exec("UPDATE posts SET Likes = Likes + $1, Dislikes = Dislikes - $2 WHERE id = $3", like, dislike, string_id_ofpost); err != nil {
		return err
	}
	return nil
}

func (s storage) UpdateDislikePost(id_of_post, like, dislike int) error {
	string_id_ofpost := strconv.Itoa(id_of_post)
	if _, err := s.db.Exec("UPDATE posts SET Likes = Likes - $1, Dislikes = Dislikes + $2 WHERE id = $3", like, dislike, string_id_ofpost); err != nil {
		return err
	}
	return nil
}

func (s storage) PostIsLike(username string, post_id int) (int, bool) {
	string_id_ofpost := strconv.Itoa(post_id)
	num := 0
	err := s.db.QueryRow("SELECT IsLike FROM posts_is_like WHERE Username = ? AND Post_id=?", username, string_id_ofpost).Scan(&num)
	if err == sql.ErrNoRows {
		return 0, false
	}
	return num, true
}

func (s storage) InsertUserLike(like int, username string, post_id int) error {
	string_id_ofpost := strconv.Itoa(post_id)
	if _, err := s.db.Exec("INSERT INTO posts_is_like (IsLike,Username,Post_id) VALUES ($1,$2,$3)", like, username, string_id_ofpost); err != nil {
		return err
	}

	return nil
}

func (s storage) CommentIsLike(username string, post_id int, comment_id int) (int, bool) {
	string_id_ofpost := strconv.Itoa(post_id)
	string_id_ofcomment := strconv.Itoa(comment_id)
	num := 0
	err := s.db.QueryRow("SELECT IsLike FROM comments_is_like WHERE Username = ? AND Post_id=? AND Comment_id = ?", username, string_id_ofpost, string_id_ofcomment).Scan(&num)
	if err == sql.ErrNoRows {
		return 0, false
	}
	return num, true
}

func (s storage) UpdateIsLikeComment(username string, post_id, comment_id, like int) error {
	string_id_ofpost := strconv.Itoa(post_id)
	string_id_ofcomment := strconv.Itoa(comment_id)
	if _, err := s.db.Exec("UPDATE comments_is_like SET IsLike = ? WHERE Username = ? AND Post_id = ? AND Comment_id = ?", like, username, string_id_ofpost, string_id_ofcomment); err != nil {
		return err
	}
	return nil
}

func (s storage) UpdateLikeComment(comment_id, like, dislike int) error {
	string_id_ofcomment := strconv.Itoa(comment_id)
	if _, err := s.db.Exec("UPDATE comments SET Likes = Likes + $1, Dislikes = Dislikes - $2 WHERE id = $3", like, dislike, string_id_ofcomment); err != nil {
		return err
	}
	return nil
}

func (s storage) UpdateDislikeComment(id_of_comment, like, dislike int) error {
	string_id_ofcomment := strconv.Itoa(id_of_comment)
	if _, err := s.db.Exec("UPDATE comments SET Likes = Likes - $1, Dislikes = Dislikes + $2 WHERE id = $3", like, dislike, string_id_ofcomment); err != nil {
		return err
	}
	return nil
}

func (s storage) InsertUserLikeComment(like int, username string, post_id int, comment_id int) error {
	string_id_ofpost := strconv.Itoa(post_id)
	string_id_ofcomment := strconv.Itoa(comment_id)
	if _, err := s.db.Exec("INSERT INTO comments_is_like (IsLike,Username,Post_id,Comment_id) VALUES ($1,$2,$3,$4)", like, username, string_id_ofpost, string_id_ofcomment); err != nil {
		return err
	}

	return nil
}

func (s storage) GetLikedPosts(username string) ([]int, error) {
	id_posts := []int{}
	rows, err := s.db.Query("SELECT Post_id FROM posts_is_like WHERE IsLike = 1 AND Username like '%" + username + "%'")
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		id_posts = append(id_posts, id)
	}
	return id_posts, nil
}
