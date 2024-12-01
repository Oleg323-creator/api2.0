CREATE TABLE IF NOT EXISTS rates(
     id SERIAL PRIMARY KEY,
     created_at TIMESTAMP,
     updated_at TIMESTAMP,
     "from" INT,
     "to" INT,
     rate INT,
     provider VARCHAR(100)
);