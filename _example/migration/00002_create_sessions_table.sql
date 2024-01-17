-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions
(
    id      SERIAL PRIMARY KEY,
    token   VARCHAR(255) NOT NULL UNIQUE,
    user_id INT          NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
