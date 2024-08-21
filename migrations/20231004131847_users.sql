-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        login VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        registered_at timestamp NOT NULL
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd