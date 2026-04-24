package repositories

import (
	"currency-exchange/internal/domain"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
)

type ExchangeRateRepository struct {
	db *sql.DB
}

func ExchangeRateRepositoryNew(db *sql.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) AddExchangeRates(baseCurrency domain.Currency, targetCurrency domain.Currency, rate *big.Rat) (domain.ExchangeRate, error) {
	query := `INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) VALUES (?, ?, ?)`

	res, err := r.db.Exec(query, baseCurrency.ID, targetCurrency.ID, rate.FloatString(4))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.ExchangeRate{}, domain.ErrExchangeRateAlreadyExists
		}
		return domain.ExchangeRate{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return domain.ExchangeRate{
		ID:               int(id),
		BaseCurrencyId:   baseCurrency.ID,
		TargetCurrencyId: targetCurrency.ID,
		Rate:             rate,
	}, nil
}

func (r *ExchangeRateRepository) GetExchangeRates() ([]domain.ExchangeRate, error) {
	query := `SELECT * FROM exchange_rates`

	rows, err := r.db.Query(query)
	result := []domain.ExchangeRate{}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	for rows.Next() {
		var rateStr string
		var e domain.ExchangeRate

		err := rows.Scan(
			&e.ID,
			&e.BaseCurrencyId,
			&e.TargetCurrencyId,
			&rateStr,
		)
		if err != nil {
			return nil, err
		}

		rate := new(big.Rat)
		_, ok := rate.SetString(rateStr)
		if !ok {
			return []domain.ExchangeRate{}, fmt.Errorf("invalid rate value")
		}
		e.Rate = rate

		result = append(result, e)
	}

	return result, nil
}

func (r *ExchangeRateRepository) GetExchangeRate(baseCurrencyId int, targetCurrencyId int) (domain.ExchangeRate, bool, error) {
	var e domain.ExchangeRate
	query := `SELECT id, base_currency_id, target_currency_id, rate FROM exchange_rates WHERE base_currency_id = ? AND target_currency_id = ?`

	var rateStr string

	err := r.db.QueryRow(query, baseCurrencyId, targetCurrencyId).Scan(
		&e.ID,
		&e.BaseCurrencyId,
		&e.TargetCurrencyId,
		&rateStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ExchangeRate{}, false, nil
		}
		return domain.ExchangeRate{}, false, err
	}

	rate := new(big.Rat)
	_, ok := rate.SetString(rateStr)
	if !ok {
		return domain.ExchangeRate{}, false, fmt.Errorf("invalid rate value")
	}
	e.Rate = rate
	return e, true, nil
}

func (r *ExchangeRateRepository) UpdateExchangeRate(baseCurrency domain.Currency, targetCurrency domain.Currency, rate *big.Rat) (domain.ExchangeRate, error) {
	exchangeRate, found, err := r.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
	if !found {
		return domain.ExchangeRate{}, domain.ErrExchangeRateNotFound
	}
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	query := `UPDATE exchange_rates SET rate = ? WHERE id = ?`

	res, err := r.db.Exec(query, rate.FloatString(4), exchangeRate.ID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	_, _ = res.RowsAffected()

	exchangeRate.Rate = rate
	return exchangeRate, nil
}
