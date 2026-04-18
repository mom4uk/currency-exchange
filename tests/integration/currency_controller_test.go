package integration

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/test_utilities"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGetCurrencies_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := test_utilities.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected: 200, got: %d", rr.Code)
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

func TestGetCurrency_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := test_utilities.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/currency/EUR", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
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

func TestAddCurrency_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := test_utilities.SeedCurrencies(app.DB); err != nil {
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
		t.Fatalf("expected 201, got %d", rr.Code)
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
