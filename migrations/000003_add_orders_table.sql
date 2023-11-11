-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
   id VARCHAR(30) PRIMARY KEY,
   user_id VARCHAR(36) NOT NULL,
   status VARCHAR(30) NOT NULL,
   accrual DECIMAL(10, 2) NOT NULL DEFAULT 0,
   uploaded_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_id_user_id on orders (id, user_id);
CREATE INDEX IF NOT EXISTS idx_user_id on orders (user_id);
CREATE INDEX IF NOT EXISTS idx_status on orders (status);


-- +goose Down
DROP TABLE users;
DROP INDEX IF EXISTS uniq_idx_id_user_id;
DROP INDEX IF EXISTS idx_user_id;
DROP INDEX IF EXISTS idx_status;
