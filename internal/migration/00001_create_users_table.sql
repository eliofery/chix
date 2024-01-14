-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    age  SMALLINT    NOT NULL
);

INSERT INTO users (name, age) VALUES ('John', 30), ('Jane', 25), ('Bob', 40), ('Alice', 35);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
