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

func (r *CurrencyRepository) AddCurrency(c domain.Currency) error {
	query := `INSERT INTO currencies (name, code, sign) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, c.Name, c.Code, c.Sign)
	return err
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

func (r *CurrencyRepository) AddExchangeRates(e domain.AddExchangeRateRequest) (domain.ExchangeRateResponse, error) {
	baseCurrency, err := r.GetCurrencyByCode(e.BaseCurrencyCode)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	targetCurrency, err := r.GetCurrencyByCode(e.TargetCurrencyCode)
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

func (r *CurrencyRepository) GetExchangeRates() ([]domain.ExchangeRateResponse, error) {
	query := `SELECT * FROM exchange_rates`

	rows, err := r.db.Query(query)
	result := []domain.ExchangeRateResponse{}

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

		baseCurrency, err := r.GetCurrencyById(e.BaseCurrencyId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			return nil, err
		}
		targetCurrency, err := r.GetCurrencyById(e.TargetCurrencyId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			return nil, err
		}

		var res domain.ExchangeRateResponse

		res.BaseCurrency = baseCurrency
		res.TargetCurrency = targetCurrency
		res.Rate = e.Rate

		result = append(result, res)
	}

	return result, nil
}

func (r *CurrencyRepository) GetExchangeRatesByCodes(baseCurrencyCode string, targetCurrencyCode string) (domain.ExchangeRateResponse, error) {
	baseCurrency, err := r.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	targetCurrency, err := r.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	rate, err := r.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}
	return domain.ExchangeRateResponse{
		ID:             int(rate.ID),
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate.Rate,
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
