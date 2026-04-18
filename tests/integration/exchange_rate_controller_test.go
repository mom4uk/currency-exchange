package integration

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/test_utilities"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetExchangeRates_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := test_utilities.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := test_utilities.SeedExchangeRates(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/exchangeRates", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got []domain.ExchangeRateResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []domain.ExchangeRateResponse{
		{
			ID: 1,
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
			Rate: 0.99,
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}

func TestGetExchangeRate_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	if err := test_utilities.SeedCurrencies(app.DB); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	if err := test_utilities.SeedExchangeRates(app.DB); err != nil {
		t.Fatalf("failed to seed exchange rates: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/exchangeRate/USDEUR", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	var got domain.ExchangeRateResponse
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	expected := domain.ExchangeRateResponse{
		ID: 1,
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
		Rate: 0.99,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %+v\nexpected: %+v", got, expected)
	}
}
