package repositories

import (
	"database/sql"
)

type CurrencyRepository struct {
	db *sql.DB
}

type Currency struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Sign string `json:"sign"`
}

func CurrencyRepositoryNew(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r CurrencyRepository) GetCurrencies() ([]Currency, error) {
	query := "SELECT * FROM currencies"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []Currency{}

	for rows.Next() {
		var c Currency
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
