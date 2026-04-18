package test_utilities

import (
	"currency-exchange/internal/controllers"
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
	"currency-exchange/internal/server"
	"currency-exchange/internal/services"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func NewTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	if err := runMigrations(db); err != nil {
		t.Fatal(err)
	}

	return db
}

type TestApp struct {
	DB                  *sql.DB
	Server              *server.Server
	CurrencyService     *services.CurrencyService
	ExchangeRateService *services.ExchangeRateService
}

func NewTestApp(t *testing.T) *TestApp {
	db := NewTestDB(t)

	currencyRepo := repositories.CurrencyRepositoryNew(db)
	exchangeRepo := repositories.ExchangeRateRepositoryNew(db)

	currencyService := services.CurrencyServiceNew(exchangeRepo, currencyRepo)
	exchangeRateService := services.ExchangeRateServiceNew(exchangeRepo, currencyRepo)
	exchangeService := services.ExchangeServiceNew(exchangeRepo, currencyRepo)

	currencyController := controllers.NewController(currencyService)
	exchangeRateController := controllers.NewExchangeRateController(exchangeRateService)
	exchangeController := controllers.NewExchangeController(exchangeService)

	srv := server.New()
	controllers.RegisterCurrencyRoutes(srv.GetMux(), currencyController)
	controllers.RegisterExchangeRateRoutes(srv.GetMux(), exchangeRateController)
	controllers.RegisterExchangeRoutes(srv.GetMux(), exchangeController)

	return &TestApp{
		DB:                  db,
		Server:              srv,
		CurrencyService:     currencyService,
		ExchangeRateService: exchangeRateService,
	}
}

func runMigrations(db *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	base := filepath.Dir(filename)

	migrationsPath := filepath.Join(base, "../../db/migrations")

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if strings.Contains(file.Name(), ".down.") {
			continue
		}
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		fullPath := filepath.Join(migrationsPath, file.Name())

		query, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(query)); err != nil {
			return err
		}
	}

	return nil
}

func AssertCurrencies(t *testing.T, got, exp []domain.Currency) {
	t.Helper()

	if len(got) != len(exp) {
		t.Fatalf("expected %d currencies, got %d\n got: %+v", len(exp), len(got), got)
	}

	used := make([]bool, len(exp))

	for _, c := range got {
		found := false

		for i, e := range exp {
			if used[i] {
				continue
			}

			if c.Code == e.Code {
				if c.Name != e.Name || c.Sign != e.Sign {
					t.Fatalf(
						"currency mismatch for code=%s:\n got: %+v\n exp: %+v",
						c.Code, c, e,
					)
				}

				used[i] = true
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("unexpected currency: %+v", c)
		}
	}

	for i, ok := range used {
		if !ok {
			t.Fatalf("missing currency: %+v", exp[i])
		}
	}
}
