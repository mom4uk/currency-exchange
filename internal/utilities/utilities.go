package utilities

import (
	"currency-exchange/internal/domain"
	"encoding/json"
	"errors"
	"math/big"
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

func WriteError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(domain.ErrorResponse{
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
		return

	case errors.Is(err, domain.ErrCurrencyAlreadyExists):
		WriteError(w, "Такая валюта уже существует", http.StatusConflict)
		return

	case errors.Is(err, domain.ErrExchangeRateNotFound):
		WriteError(w, "Такой обменный курс не найден", http.StatusNotFound)
		return

	case errors.Is(err, domain.ErrExchangeRateAlreadyExists):
		WriteError(w, "Такой обменный курс уже существует", http.StatusConflict)
		return

	case errors.Is(err, domain.ErrAbsenceOfExchangeRateField):
		WriteError(w, "Отстутствует одно из обязательных полей: baseCurrencyCode, targetCurrencyCode, rate", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrAbsenceOfExchangeRateFieldForUpdate):
		WriteError(w, "Отстутствует обязательное поле: rate", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrInvalidCurrencyField):
		WriteError(w, "Некорректные значения в полях", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrInvalidCurrencySign):
		WriteError(w, "Значение в поле sign не должно быть длинее 3 символов", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrAmountFormatIncorrect):
		WriteError(w, "Некорректное значение amount", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrMissingFromCurrency):
		WriteError(w, "Значение from не передано", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrMissingToCurrency):
		WriteError(w, "Значение to не передано", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrMissingAmount):
		WriteError(w, "Значение amount не передано", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrAmountConvertation):
		WriteError(w, "Ошибка в конвертации amount", http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrRateConvertaion):
		WriteError(w, "Ошибка в конвертации rate", http.StatusBadRequest)
		return

	default:
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RatToFloat(r *big.Rat) float64 {
	f, _ := r.Float64()
	return f
}
