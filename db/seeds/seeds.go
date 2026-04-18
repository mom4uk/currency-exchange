package seeds

import "database/sql"

func SeedRubCurrency(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO currencies (code, name, sign) VALUES
		('RUB', 'Russian Ruble', '₽');
	`)
	if err != nil {
		return err
	}

	return nil
}

func SeedExchangeEurToUsd(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate)
		VALUES (2, 1, 2);
	`)
	if err != nil {
		return err
	}

	return nil
}

func SeedExchangeCrossViaUsd(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) VALUES
		(1, 2, 0.5),
		(1, 3, 100);
	`)
	if err != nil {
		return err
	}

	return nil
}

func SeedCurrencies(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO currencies (code, name, sign) VALUES
		('USD', 'United States dollar', '$'),
		('EUR', 'Euro', '€');
	`)
	if err != nil {
		return err
	}

	return nil
}

func SeedExchangeUsdToEur(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) 
		VALUES (1, 2, 0.99);
	`)
	if err != nil {
		return err
	}

	return nil
}
