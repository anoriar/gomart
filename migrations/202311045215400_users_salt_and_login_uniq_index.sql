-- +goose Up
CREATE UNIQUE INDEX uniq_ids_users_login ON users (login);

ALTER TABLE users ADD COLUMN salt VARCHAR(255);



-- +goose Down
DROP INDEX IF EXISTS uniq_ids_users_login;
ALTER TABLE users DROP COLUMN salt;

