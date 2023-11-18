-- +goose Up
CREATE TABLE IF NOT EXISTS withdrawals (
   id UUID PRIMARY KEY,
   user_id VARCHAR(36) NOT NULL,
   order_id VARCHAR(30) NOT NULL,
   sum DECIMAL(10, 2) NOT NULL DEFAULT 0,
   processed_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_withdrawals_user_id on withdrawals (user_id);
CREATE INDEX IF NOT EXISTS idx_withdrawals_user_id_sum on withdrawals (user_id, sum);


-- +goose Down
DROP TABLE withdrawals;
DROP INDEX IF EXISTS idx_withdrawals_user_id;
