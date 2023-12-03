-- +goose Up
CREATE TABLE IF NOT EXISTS users (
   id VARCHAR(36) PRIMARY KEY,
   login VARCHAR(255) NOT NULL,
   password VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
