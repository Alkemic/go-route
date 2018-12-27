package middleware

import "net/http"

func Noop(f http.HandlerFunc) http.HandlerFunc {
	return f
}
