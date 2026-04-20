package integration

import (
	"currency-exchange/db/seeds"
	"currency-exchange/internal/domain"
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

	var got domain.CurencyExchange
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := domain.CurencyExchange{
		BaseCurrency: domain.Currency{
			ID:   1,
			Code: "USD",
			Name: "United States dollar",
			Sign: "$",
		},
		TargetCurrency: domain.Currency{
			ID:   2,
			Code: "EUR",
			Name: "Euro",
			Sign: "€",
		},
		Rate:            0.99,
		Amount:          10,
		ConvertedAmount: 9.9,
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

	var got domain.CurencyExchange
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedRate := 2.0

	if got.Rate != expectedRate {
		t.Fatalf("expected rate %v, got %v", expectedRate, got.Rate)
	}

	if got.ConvertedAmount != 20 {
		t.Fatalf("expected 20, got %v", got.ConvertedAmount)
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

	var got domain.CurencyExchange
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedRate := 200.0

	if got.Rate != expectedRate {
		t.Fatalf("expected rate %v, got %v", expectedRate, got.Rate)
	}

	if got.ConvertedAmount != 2000 {
		t.Fatalf("expected 2000, got %v", got.ConvertedAmount)
	}
}
