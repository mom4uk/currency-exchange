package utilities

import (
	"currency-exchange/internal/domain"
	"net/http"
	"strings"
)

func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func GetLastPathSegment(url string) (string, error) {
	index := strings.LastIndex(url, "/")
	if url[index+1:] == "" {
		return "", domain.ErrAbsenceOfCode
	}
	return url[index+1:], nil
}

func GetCurrencyCodes(url string) (string, string, error) {
	codes, err := GetLastPathSegment(url)
	if err != nil {
		return "", "", domain.ErrAbsenceOfCode
	}

	if len(codes) != 6 {
		return "", "", domain.ErrIncorrectLengthOfCode
	}

	return codes[:3], codes[3:], nil
}
