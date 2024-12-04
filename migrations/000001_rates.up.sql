CREATE TABLE IF NOT EXISTS rates(
     id SERIAL PRIMARY KEY,
     created_at TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     from_currency VARCHAR(20),
     to_currency VARCHAR(20),
     rate DOUBLE PRECISION,
     provider VARCHAR(100),
     CONSTRAINT unique_cols UNIQUE (provider, from_currency, to_currency)
);