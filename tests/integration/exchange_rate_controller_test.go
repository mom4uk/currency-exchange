package integration

import (
	"currency-exchange/db/seeds"
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/test_utilities"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// Success for GET /exchangeRates
func TestGetExchangeRates_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/exchangeRates", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got []dto.ExchangeRateResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []dto.ExchangeRateResponse{
		{
			ID: 1,
			BaseCurrency: dto.CurrencyResponse{
				ID:   "1",
				Code: "USD",
				Name: "United States dollar",
				Sign: "$",
			},
			TargetCurrency: dto.CurrencyResponse{
				ID:   "2",
				Code: "EUR",
				Name: "Euro",
				Sign: "€",
			},
			Rate: 0.9900,
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestGetExchangeRates_responseNumericFields(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/exchangeRates", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var raw []map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for i, item := range raw {
		rate, ok := item["rate"]
		if !ok {
			t.Fatalf("item %d: rate missing", i)
		}

		if _, ok := rate.(float64); !ok {
			t.Fatalf("item %d: rate expected number, got %T (%v)", i, rate, rate)
		}
	}
}

// Success for GET /exchangeRate/{code}
func TestGetExchangeRate_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/exchangeRate/USDEUR", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got dto.ExchangeRateResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := dto.ExchangeRateResponse{
		ID: 1,
		BaseCurrency: dto.CurrencyResponse{
			ID:   "1",
			Code: "USD",
			Name: "United States dollar",
			Sign: "$",
		},
		TargetCurrency: dto.CurrencyResponse{
			ID:   "2",
			Code: "EUR",
			Name: "Euro",
			Sign: "€",
		},
		Rate: 0.9900,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

// Errors for GET /exchangeRate/{code}
func TestGetExchangeRate_error_absenceOfCodes(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/exchangeRate/", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Вы не передали код валюты",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

func TestGetExchangeRate_error_exchangeRateNotFound(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/exchangeRate/EURUSD", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Такой обменный курс не найден",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

// Success for POST /exchangeRates
func TestAddExchangeRate_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	form := url.Values{}
	form.Add("baseCurrencyCode", "USD")
	form.Add("targetCurrencyCode", "EUR")
	form.Add("rate", "0.9900")

	req := httptest.NewRequest(
		http.MethodPost,
		"/exchangeRates",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got dto.ExchangeRateResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := dto.ExchangeRateResponse{
		ID: 1,
		BaseCurrency: dto.CurrencyResponse{
			ID:   "1",
			Code: "USD",
			Name: "United States dollar",
			Sign: "$",
		},
		TargetCurrency: dto.CurrencyResponse{
			ID:   "2",
			Code: "EUR",
			Name: "Euro",
			Sign: "€",
		},
		Rate: 0.9900,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

// Errors for POST /exchangeRates
func TestPostExchangeRate_error_currencyNotFound(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedRubCurrency(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	form := url.Values{}

	form.Add("baseCurrencyCode", "USD")
	form.Add("targetCurrencyCode", "RUB")
	form.Add("rate", "0.9900")

	req := httptest.NewRequest(
		http.MethodPost,
		"/exchangeRates",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Такая валюта не найдена",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

func TestPostExchangeRate_error_absenceOfFields(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	form := url.Values{}
	form.Add("baseCurrencyCode", "USD")
	form.Add("rate", "0.99")

	req := httptest.NewRequest(
		http.MethodPost,
		"/exchangeRates",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Отстутствует одно из обязательных полей: baseCurrencyCode, targetCurrencyCode, rate",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

func TestPostExchangeRate_error_exchangeRateAlreadyExists(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	form := url.Values{}
	form.Add("baseCurrencyCode", "USD")
	form.Add("targetCurrencyCode", "EUR")
	form.Add("rate", "0.9900")

	req := httptest.NewRequest(
		http.MethodPost,
		"/exchangeRates",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Такой обменный курс уже существует",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

// Success for PATCH /exchangeRates/{code}
func TestUpdateExchangeRate_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	form := url.Values{}
	form.Add("rate", "0.98")

	req := httptest.NewRequest(
		http.MethodPatch,
		"/exchangeRate/USDEUR",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got dto.ExchangeRateResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := dto.ExchangeRateResponse{
		ID: 1,
		BaseCurrency: dto.CurrencyResponse{
			ID:   "1",
			Code: "USD",
			Name: "United States dollar",
			Sign: "$",
		},
		TargetCurrency: dto.CurrencyResponse{
			ID:   "2",
			Code: "EUR",
			Name: "Euro",
			Sign: "€",
		},
		Rate: 0.9800,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

// Errors for PATCH /exchangeRates/{code}

func TestPatchExchangeRate_error_absenceOfFields(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	form := url.Values{}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/exchangeRate/USDEUR",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Отстутствует обязательное поле: rate",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

func TestPatchExchangeRate_error_currencyNotFound(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	form := url.Values{}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/exchangeRate/USDEUR",
		strings.NewReader(form.Encode()),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 404, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v\nbody: %s", err, rr.Body.String())
	}

	expected := domain.ErrorResponse{
		Message: "Отстутствует обязательное поле: rate",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}
