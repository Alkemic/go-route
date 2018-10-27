package middleware

import (
	"net/http"
)

func AllowedMethods(allowedMethods []string) func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			for _, method := range allowedMethods {
				if req.Method == method {
					f(rw, req)
					return
				}
			}

			http.Error(rw, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
