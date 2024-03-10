package middleware

import (
	"log"
	"net/http"
)

func LogRoute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("(%s) %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}
