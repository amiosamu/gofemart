-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    withdrawals (
        order_id VARCHAR(255) NOT NULL UNIQUE,
        bonuses numeric,
        uploaded_at TIMESTAMPTZ NOT NULL,
        user_id integer REFERENCES users (id)
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS withdrawals;

-- +goose StatementEnd