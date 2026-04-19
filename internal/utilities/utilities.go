package utilities

import (
	"currency-exchange/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

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

func WriteError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	})
}

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCurrencyNotFound):
		WriteError(w, "Такая валюта не найдена", http.StatusNotFound)
		return

	case errors.Is(err, domain.ErrAbsenceOfCode):
		WriteError(w, "Вы не передали код валюты", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrIncorrectLengthOfCode):
		WriteError(w, "Неверный формат кода", http.StatusBadRequest)
		return

	default:
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
