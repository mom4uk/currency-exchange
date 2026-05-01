package apierrors

import (
	"currency-exchange/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

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
		WriteError(w, "internal error", http.StatusInternalServerError)
		return
	}
}
