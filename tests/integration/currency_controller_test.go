package integration

import (
	"currency-exchange/db/seeds"
	"currency-exchange/internal/domain"
	"currency-exchange/internal/test_utilities"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Success for GET /currencies
func TestGetCurrencies_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var resp []domain.Currency
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []domain.Currency{
		{Code: "USD", Name: "United States dollar", Sign: "$"},
		{Code: "EUR", Name: "Euro", Sign: "€"},
	}

	test_utilities.AssertCurrencies(t, resp, expected)
}

// Success for GET /currency/{code}
func TestGetCurrency_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/currency/EUR", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.Currency
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := domain.Currency{
		ID:   2,
		Code: "EUR",
		Name: "Euro",
		Sign: "€",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

// Errors for GET /currency/{code}

func TestGetCurrency_error_absenceOfCurrencyCode(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/currency/", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := domain.ErrorResponse{
		Message: "Вы не передали код валюты",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

func TestGetCurrency_error_currencyNotFound(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/currency/EUR", nil)
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

// Success for POST /currency/{code}
func TestAddCurrency_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/currencies",
		strings.NewReader("name=Russian+Ruble&code=RUB&sign=₽"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.Currency
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := domain.Currency{
		ID:   3,
		Code: "RUB",
		Name: "Russian Ruble",
		Sign: "₽",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

// Errors for POST /currency/{code}
func TestAddCurrency_error_abscenceOfFormFields(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(
		http.MethodPost,
		"/currencies",
		strings.NewReader("name=Russian+Ruble"),
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
		Message: "Отстутствует одно из обязательных полей: name, code, sign",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}

func TestAddCurrency_error_currencyAlreadyExists(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/currencies",
		strings.NewReader("name=Euro&code=EUR&sign=€"),
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
		Message: "Такая валюта уже существует",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v, expected: %+v", got, expected)
	}
}
