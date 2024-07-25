-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox (
    id           BIGSERIAL PRIMARY KEY,
    topic        TEXT NOT NULL,
    event        BYTEA NOT NULL,
    headers      BYTEA NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    completed_at TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd
