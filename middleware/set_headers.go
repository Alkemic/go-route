package middleware

import (
	"net/http"

	"github.com/Alkemic/go-route"
)

func SetHeaders(headers map[string]string) func(f route.HandlerFunc) route.HandlerFunc {
	return func(f route.HandlerFunc) route.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			for k, v := range headers {
				rw.Header().Set(k, v)
			}

			f(rw, req)
		}
	}
}
