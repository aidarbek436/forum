package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aidarbek436/forum/internal/repository"
)

func (h *handler) RegistrationHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		err := ErrorPage(w, 404, "404 Page not Found:/register")
		if err != nil {
			fmt.Println(err)
		}
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
		absPath, err := filepath.Abs("front/register.html")
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:FilePath")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		temp, err := template.ParseFiles(absPath)
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:Parse template")
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		err = temp.Execute(w, nil)
		if err != nil {
			err := ErrorPage(w, 500, "Internal server error:template execute")
			if err != nil {
				fmt.Println(err)
			}
			return
		}

	case http.MethodPost:
		_, err1 := r.Cookie("forum")
		if err1 == nil {
			err := ErrorPage(w, 400, "Bad request:session already exists")
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		username := r.FormValue("username")
		usernames := strings.ReplaceAll(username, " ", "")
		if usernames == "" {
			err := ErrorPage(w, 400, "400 Bad Request: nil username value")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if HtmlInjectionCheck(username) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		email := r.FormValue("email")
		emails := strings.ReplaceAll(email, " ", "")
		if emails == "" {
			err := ErrorPage(w, 400, "400 Bad Request: nil email value")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if HtmlInjectionCheck(email) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		password := r.FormValue("password")
		passwords := strings.ReplaceAll(password, " ", "")
		if passwords == "" {
			err := ErrorPage(w, 400, "400 Bad Request: nil password value")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if HtmlInjectionCheck(password) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		confPass := r.FormValue("confirm_password")
		confPasss := strings.ReplaceAll(confPass, " ", "")
		if confPasss == "" {
			err := ErrorPage(w, 400, "400 Bad Request: nil password value")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if HtmlInjectionCheck(confPass) == false {
			err := ErrorPage(w, 400, "400 Bad Request: html injection")
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		if !h.storage.UserIsExist(username) {
			err := ErrorPage(w, 400, "400 Bad Request:User is exists")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("User is exist")
			return
		}
		if !h.storage.EmailIsExist(email) {
			err := ErrorPage(w, 400, "400 Bad Request: Email is exists")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Email is exist")
			return
		}
		err := emailIsvalid(email)
		if err != nil {
			err := ErrorPage(w, 400, "400 Bad Request: email is invalid")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err = passIsValid(password)
		if err != nil {
			err := ErrorPage(w, 400, "400 Bad Request:password is invalid")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if password != confPass {
			err := ErrorPage(w, 400, "400 Bad Request:passwords do not match")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		var userInput repository.User // todo 1)handle r body 2) handle http
		userInput.Username = username
		userInput.Email = email
		userInput.Password = password
		if err = h.storage.PostUser(userInput); err != nil {
			err := ErrorPage(w, 500, "Internal server error:Post user")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	default:
		err := ErrorPage(w, 405, "Method not Allowed:should be Get")
		if err != nil {
			fmt.Println(err)
		}
		return

	}
}

func passIsValid(pass string) error {
	if len(pass) < 8 {
		return errors.New("Password should be of 8 characters long")
	}
	done, err := regexp.MatchString("([a-z])+", pass)
	if err != nil {
		return err
	}
	if !done {
		return errors.New("Password should contain atleast one lower case character")
	}
	done, err = regexp.MatchString("([0-9])+", pass)
	if err != nil {
		return err
	}
	if !done {
		return errors.New("Password should contain atleast one digit")
	}
	return nil
}

func emailIsvalid(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
