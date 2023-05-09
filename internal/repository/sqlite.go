package repository

import (
	"database/sql"

	"github.com/aidarbek436/forum/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDb(cfg config.Database) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBname)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS "users" (
		"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"Username"	TEXT,
		"Email"	TEXT,
		"Password"	TEXT
		
	);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	query1 := `
	CREATE TABLE IF NOT EXISTS "posts" (
		"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"Title"	TEXT,
		"Content"	TEXT,
		"Likes" INTEGER,
		"Dislikes" INTEGER,
		"Author"	TEXT,
		"IdOfAuthor" INTEGER,
		"C_sport"    INTEGER,
		"C_history"  INTEGER,
		"C_politics" INTEGER,
		"C_science"  INTEGER,
		"C_art"      INTEGER
	);`

	_, err = db.Exec(query1)
	if err != nil {
		return nil, err
	}
	query2 := `
	CREATE TABLE IF NOT EXISTS "comments" (
		"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"Content"	TEXT,
		"Likes" INTEGER,
		"Dislikes" INTEGER,
		"Author"	TEXT,
		"IdOfPost" INTEGER
	);`
	_, err = db.Exec(query2)
	if err != nil {
		return nil, err
	}
	query3 := `
	CREATE TABLE IF NOT EXISTS "posts_is_like" (
		"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"IsLike" INTEGER,
		"Username" TEXT,
		"Post_id" INTEGER
	);`
	_, err = db.Exec(query3)
	if err != nil {
		return nil, err
	}
	query4 := `
	CREATE TABLE IF NOT EXISTS "comments_is_like" (
		"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"IsLike" INTEGER,
		"Username" TEXT,
		"Post_id" INTEGER,
		"Comment_id" INTEGER
	);`
	_, err = db.Exec(query4)
	if err != nil {
		return nil, err
	}
	// statement.Exec()
	return db, nil
}
