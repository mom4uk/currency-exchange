package integration

import (
	"currency-exchange/db/seeds"
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/test_utilities"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestExchange_success_directRate(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=USD&to=EUR&amount=10",
		nil,
	)

	rr := httptest.NewRecorder()
	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got dto.CurencyExchangeResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := dto.CurencyExchangeResponse{
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
		Rate:            0.99,
		Amount:          10.00,
		ConvertedAmount: 9.90,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestExchange_success_reverseRate(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeEurToUsd(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=EUR&to=USD&amount=10",
		nil,
	)

	rr := httptest.NewRecorder()
	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got dto.CurencyExchangeResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedRate := 2.00

	if got.Rate != expectedRate {
		t.Fatalf("expected rate %v, got %v", expectedRate, got.Rate)
	}

	if got.ConvertedAmount != 20.00 {
		t.Fatalf("expected 20.00, got %v", got.ConvertedAmount)
	}
}

func TestExchange_success_viaUSD(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedRubCurrency(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeCrossViaUsd(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=EUR&to=RUB&amount=10",
		nil,
	)

	rr := httptest.NewRecorder()
	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got dto.CurencyExchangeResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedRate := 200.00

	if got.Rate != expectedRate {
		t.Fatalf("expected rate %v, got %v", expectedRate, got.Rate)
	}

	if got.ConvertedAmount != 2000.00 {
		t.Fatalf("expected 2000.00, got %v", got.ConvertedAmount)
	}
}

func TestExchange_error_absenceOfAmountField(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=USD&to=EUR",
		nil,
	)

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
		Message: "Значение amount не передано",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestExchange_error_absenceOfFromField(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?to=EUR&amount=10",
		nil,
	)

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
		Message: "Значение from не передано",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestExchange_error_absenceOfToField(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=USD&amount=10",
		nil,
	)

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
		Message: "Значение to не передано",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestExchange_error_incorrectAmountValue(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=USD&to=EUR&amount=lol",
		nil,
	)

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
		Message: "Некорректное значение amount",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestExchange_success_responseNumericFields(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := seeds.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := seeds.SeedExchangeUsdToEur(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/exchange?from=USD&to=EUR&amount=10",
		nil,
	)

	rr := httptest.NewRecorder()
	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var raw map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	rate, ok := raw["rate"].(float64)
	if !ok {
		t.Fatalf("expected rate to be number, got %T (%v)", raw["rate"], raw["rate"])
	}

	amount, ok := raw["amount"].(float64)
	if !ok {
		t.Fatalf("expected amount to be number, got %T (%v)", raw["amount"], raw["amount"])
	}

	convertedAmount, ok := raw["convertedAmount"].(float64)
	if !ok {
		t.Fatalf(
			"expected convertedAmount to be number, got %T (%v)",
			raw["convertedAmount"],
			raw["convertedAmount"],
		)
	}

	if rate != 0.99 {
		t.Fatalf("expected rate 0.99, got %v", rate)
	}

	if amount != 10 {
		t.Fatalf("expected amount 10, got %v", amount)
	}

	if convertedAmount != 9.9 {
		t.Fatalf("expected convertedAmount 9.9, got %v", convertedAmount)
	}
}
