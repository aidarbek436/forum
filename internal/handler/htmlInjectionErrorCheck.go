package handler

import (
	"html"
)

func HtmlInjectionCheck(input string) bool {
	safeInput := html.EscapeString(input)
	if safeInput != input {
		return false
	}
	return true
}
