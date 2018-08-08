package middleware

import (
	"net/http"

	"github.com/Alkemic/go-route"
)

func AllowedMethods(allowedMethods []string) func(f route.HandlerFunc) route.HandlerFunc {
	return func(f route.HandlerFunc) route.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request, params map[string]string) {
			for _, method := range allowedMethods {
				if req.Method == method {
					f(rw, req, params)
					return
				}
			}

			http.Error(rw, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
