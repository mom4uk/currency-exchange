package repositories

import (
	"currency-exchange/internal/domain"
	"database/sql"
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

func (r *CurrencyRepository) AddCurrency(c domain.Currency) (domain.Currency, error) {
	query := `INSERT INTO currencies (name, code, sign) VALUES (?, ?, ?)`

	res, err := r.db.Exec(query, c.Name, c.Code, c.Sign)
	if err != nil {
		return domain.Currency{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Currency{}, err
	}
	c.ID = int(id)
	return c, err
}

func (r *CurrencyRepository) GetCurrencyByCode(code string) (domain.Currency, error) {
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

func (r *CurrencyRepository) GetCurrencyById(id int) (domain.Currency, error) {
	query := `SELECT id, code, name, sign FROM currencies WHERE id = ?`

	var c domain.Currency

	err := r.db.QueryRow(query, id).Scan(
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

func (r *CurrencyRepository) AddExchangeRates(e domain.AddExchangeRateRequest) (domain.ExchangeRate, error) {
	baseCurrency, err := r.GetCurrencyByCode(e.BaseCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	targetCurrency, err := r.GetCurrencyByCode(e.TargetCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	query := `INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) VALUES (?, ?, ?)`

	res, err := r.db.Exec(query, baseCurrency.ID, targetCurrency.ID, e.Rate)
	if err != nil {
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
		Rate:             e.Rate,
	}, nil
}

func (r *CurrencyRepository) GetExchangeRates() ([]domain.ExchangeRate, error) {
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
		var e domain.ExchangeRate

		err := rows.Scan(
			&e.ID,
			&e.BaseCurrencyId,
			&e.TargetCurrencyId,
			&e.Rate,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, nil
}

func (r *CurrencyRepository) GetExchangeRatesByCodes(baseCurrencyCode string, targetCurrencyCode string) (domain.ExchangeRate, error) {
	baseCurrency, err := r.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	targetCurrency, err := r.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	rate, err := r.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}
	return domain.ExchangeRate{
		ID:               int(rate.ID),
		BaseCurrencyId:   baseCurrency.ID,
		TargetCurrencyId: targetCurrency.ID,
		Rate:             rate.Rate,
	}, nil
}

func (r *CurrencyRepository) GetExchangeRate(baseCurrencyId int, targetCurrencyId int) (domain.ExchangeRate, error) {
	var e domain.ExchangeRate
	query := `SELECT id, base_currency_id, target_currency_id, rate FROM exchange_rates WHERE base_currency_id = ? AND target_currency_id = ?`

	err := r.db.QueryRow(query, baseCurrencyId, targetCurrencyId).Scan(
		&e.ID,
		&e.BaseCurrencyId,
		&e.TargetCurrencyId,
		&e.Rate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ExchangeRate{}, sql.ErrNoRows
		}
		return domain.ExchangeRate{}, err
	}
	return e, nil
}

func (r *CurrencyRepository) UpdateExchangeRate(baseCurrencyCode string, targetCurrencyCode string, rate float64) (domain.ExchangeRate, error) {
	exchangeRate, err := r.GetExchangeRatesByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	query := `UPDATE exchange_rates SET rate = ? WHERE id = ?`

	res, err := r.db.Exec(query, rate, exchangeRate.ID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	_, _ = res.RowsAffected()

	exchangeRate.Rate = rate
	return exchangeRate, nil
}
