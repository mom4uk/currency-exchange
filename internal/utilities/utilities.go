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
		WriteError(w, "Отсутствие или неверный формат кодов валют", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrAbsenceOfCurrencyField):
		WriteError(w, "Отстутствует одно из обязательных полей: name, code, sign", http.StatusBadRequest)

	case errors.Is(err, domain.ErrCurrencyAlreadyExists):
		WriteError(w, "Такая валюта уже существует", http.StatusConflict)

	case errors.Is(err, domain.ErrExchangeRateNotFound):
		WriteError(w, "Такой обменный курс не найден", http.StatusNotFound)

	case errors.Is(err, domain.ErrExchangeRateAlreadyExists):
		WriteError(w, "Такой обменный курс уже существует", http.StatusConflict)

	case errors.Is(err, domain.ErrAbsenceOfExchangeRateField):
		WriteError(w, "Отстутствует одно из обязательных полей: baseCurrencyCode, targetCurrencyCode, rate", http.StatusBadRequest)

	case errors.Is(err, domain.ErrAbsenceOfExchangeRateFieldForUpdate):
		WriteError(w, "Отстутствует обязательное поле: rate", http.StatusBadRequest)

	case errors.Is(err, domain.ErrInvalidCurrencyField):
		WriteError(w, "Некорректные значения в полях", http.StatusBadRequest)

	case errors.Is(err, domain.ErrInvalidCurrencySign):
		WriteError(w, "Значение в поле sign не должно быть длинее 3 символов", http.StatusBadRequest)

	default:
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
