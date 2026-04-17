package integration

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/test_utilities"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrencies_success(t *testing.T) {
	app := test_utilities.NewTestApp(t)

	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	rr := httptest.NewRecorder()

	app.Server.GetMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp []domain.Currency
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := []domain.Currency{
		{ID: 1, Code: "USD", Name: "United States dollar", Sign: "$"},
		{ID: 2, Code: "EUR", Name: "Euro", Sign: "€"},
	}

	test_utilities.AssertCurrencies(t, resp, expected)
}
