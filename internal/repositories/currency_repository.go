package repositories

import "database/sql"

type CurrencyRepository struct {
	db *sql.DB
}

func CurrencyRepositoryNew(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}
func (r CurrencyRepository) getCurrencies() map[string]string {

}
