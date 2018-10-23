package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func TimeTrack(logger *log.Logger) func(f http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func(start time.Time, name string, w http.ResponseWriter) {
				logger.Printf("%s took %s", name, time.Since(start))
			}(time.Now(), fmt.Sprintf("%s %s", r.Method, r.RequestURI), w)

			f(w, r)
		}
	}
}
