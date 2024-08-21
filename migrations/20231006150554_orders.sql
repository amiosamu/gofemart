-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    orders (
        order_id VARCHAR(255) NOT NULL UNIQUE,
        status VARCHAR(255) NOT NULL,
        uploaded_at TIMESTAMPTZ NOT NULL,
        bonuses numeric,
        user_id integer REFERENCES users (id)
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS orders;

-- +goose StatementEnd