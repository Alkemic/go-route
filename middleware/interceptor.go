package middleware

import (
	"log"
	"net/http"
)

func internalServerError(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

func panicDefer(rw http.ResponseWriter, req *http.Request, logger *log.Logger) {
	if r := recover(); r != nil {
		if logger != nil {
			logger.Printf("Panic occured:\n%s\n", r)
			logger.Println("Panic end.")
		} else {
			log.Printf("Panic occured:\n%s\n", r)
			log.Println("Panic end.")
		}
		internalServerError(rw, req)
	}
}

func PanicInterceptor(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer panicDefer(rw, req, nil)

		f(rw, req)
	}
}

func PanicInterceptorWithLogger(logger *log.Logger) func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			defer panicDefer(rw, req, logger)

			f(rw, req)
		}
	}
}
