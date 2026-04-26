package repositories

import (
	"currency-exchange/internal/domain"
	"database/sql"
	"errors"
	"log"
	"strings"
)

type CurrencyRepository struct {
	db *sql.DB
}

func CurrencyRepositoryNew(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r CurrencyRepository) GetCurrencies() ([]domain.Currency, error) {
	query := "SELECT id, name, code, sign FROM currencies"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows close error: %v", err)
		}
	}()

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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *CurrencyRepository) AddCurrency(c domain.Currency) (domain.Currency, error) {
	query := `INSERT INTO currencies (name, code, sign) VALUES (?, ?, ?)`

	res, err := r.db.Exec(query, c.Name, c.Code, c.Sign)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.Currency{}, domain.ErrCurrencyAlreadyExists
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Currency{}, domain.ErrCurrencyNotFound
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
			return domain.Currency{}, domain.ErrCurrencyNotFound
		}
		return domain.Currency{}, err
	}

	return c, nil
}
