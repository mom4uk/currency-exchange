package utilities

import (
	"net/http"
	"strings"
)

func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GetCurrencyCode(url string) string {
	index := strings.LastIndex(url, "/")
	return url[index+1:]
}
