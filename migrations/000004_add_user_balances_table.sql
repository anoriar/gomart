-- +goose Up
CREATE TABLE IF NOT EXISTS balances (
   id UUID PRIMARY KEY,
   user_id VARCHAR(36) NOT NULL,
   balance DECIMAL(10, 2) NOT NULL DEFAULT 0,
   withdrawal DECIMAL(10, 2) NOT NULL DEFAULT 0,
   updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_balances_user_id on balances (user_id);


-- +goose Down
DROP TABLE balances;
DROP INDEX IF EXISTS uniq_idx_balances_user_id;
