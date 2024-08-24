CREATE TABLE Account (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL ,
    bank_name VARCHAR(100) NOT NULL,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00
);

CREATE TABLE Category (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE Transaction (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    account_id SERIAL REFERENCES Account(id),
    category_id SERIAL REFERENCES Category(id),
    amount DECIMAL(15, 2) NOT NULL,
    description TEXT,
    transaction_type VARCHAR(10) NOT NULL, -- 'INCOME' or 'EXPENSE'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_account_username ON Account(username);
CREATE INDEX idx_category_username ON Category(username);
CREATE INDEX idx_transaction_username ON Transaction(username);