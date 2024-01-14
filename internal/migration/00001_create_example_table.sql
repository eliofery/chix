-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS example
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS example;
-- +goose StatementEnd
