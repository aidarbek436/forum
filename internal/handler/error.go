package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func ErrorPage(w http.ResponseWriter, statusCode int, statusString string) error { // DONE
	w.WriteHeader(statusCode)
	absPath, err := filepath.Abs("front/error.html")
	if err != nil {
		http.Error(w, "File not found: front/error.html", http.StatusInternalServerError)
		return err
	}
	temp, err := template.ParseFiles(absPath)
	if err != nil {
		http.Error(w, "File not found: front/error.html", http.StatusInternalServerError)
		return err
	}
	err = temp.Execute(w, statusString)
	if err != nil {
		http.Error(w, statusString, statusCode)
		return err
	}
	return nil
}
