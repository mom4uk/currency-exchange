package repositories

import (
	"currency-exchange/internal/domain"
	"database/sql"
	"fmt"
)

type CurrencyRepository struct {
	db *sql.DB
}

func CurrencyRepositoryNew(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r CurrencyRepository) GetCurrencies() ([]domain.Currency, error) {
	query := "SELECT * FROM currencies"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []domain.Currency{}

	for rows.Next() {
		var c domain.Currency
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Code,
			&c.Sign,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, nil
}

func (r *CurrencyRepository) AddCurrency(c domain.Currency) error {
	query := `INSERT INTO currencies (name, code, sign) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, c.Name, c.Code, c.Sign)
	return err
}

func (r *CurrencyRepository) GetCurrency(code string) (domain.Currency, error) {
	query := `SELECT id, code, name, sign FROM currencies WHERE code = ?`

	var c domain.Currency

	err := r.db.QueryRow(query, code).Scan(
		&c.ID,
		&c.Code,
		&c.Name,
		&c.Sign,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Currency{}, sql.ErrNoRows
		}
		return domain.Currency{}, err
	}

	return c, nil
}

func (r *CurrencyRepository) AddExchangeRates(e domain.AddExchangeRateRequest) (domain.ExchangeRateResponse, error) {
	fmt.Print(e)
	baseCurrency, err := r.GetCurrency(e.BaseCurrencyCode)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	targetCurrency, err := r.GetCurrency(e.TargetCurrencyCode)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	query := `INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) VALUES (?, ?, ?)`

	res, err := r.db.Exec(query, baseCurrency.ID, targetCurrency.ID, e.Rate)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	return domain.ExchangeRateResponse{
		ID:             int(id),
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           e.Rate,
	}, nil
}
