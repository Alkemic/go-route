package middleware

import (
	"net/http"
)

func SetHeaders(headers map[string]string) func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			for k, v := range headers {
				rw.Header().Set(k, v)
			}

			f(rw, req)
		}
	}
}
