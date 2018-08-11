package middleware

import (
	"log"
	"net/http"

	"github.com/Alkemic/go-route"
)

func panicDefer(rw http.ResponseWriter, req *http.Request, logger *log.Logger) {
	if r := recover(); r != nil {
		if logger != nil {
			logger.Printf("Panic occured:\n%s\n", r)
			logger.Println("Panic end.")
		} else {
			log.Printf("Panic occured:\n%s\n", r)
			log.Println("Panic end.")
		}
		route.InternalServerError(rw, req)
	}
}

func PanicInterceptor(f route.HandlerFunc) route.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request, params map[string]string) {
		defer panicDefer(rw, req, nil)

		f(rw, req, params)
	}
}

func PanicInterceptorWithLogger(logger *log.Logger) func(f route.HandlerFunc) route.HandlerFunc {
	return func(f route.HandlerFunc) route.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request, params map[string]string) {
			defer panicDefer(rw, req, logger)

			f(rw, req, params)
		}
	}
}
