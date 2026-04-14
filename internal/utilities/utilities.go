package utilities

import (
	"fmt"
	"net/http"
	"strings"
)

func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GetLastPathSegment(url string) string {
	index := strings.LastIndex(url, "/")
	return url[index+1:]
}

func GetCurrencyCodes(url string) (string, string, error) {
	codes := GetLastPathSegment(url)

	if len(codes) != 6 {
		return "", "", fmt.Errorf("invalid currency format")
	}

	return codes[:3], codes[3:], nil
}
