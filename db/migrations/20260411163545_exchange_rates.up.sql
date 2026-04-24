CREATE TABLE exchange_rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    base_currency_id INTEGER NOT NULL,
    target_currency_id INTEGER NOT NULL,
    rate TEXT NOT NULL,

    FOREIGN KEY (base_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (target_currency_id) REFERENCES currencies(id)

    UNIQUE (base_currency_id, target_currency_id)
);