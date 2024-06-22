-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stock (
    sku         BIGINT NOT NULL,
    total_count INT NOT NULL DEFAULT 0,
    reserved    INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- TEST DATA
INSERT INTO stock (sku, total_count, reserved)
VALUES
    (1, 7, 3),
    (2, 5, 2),
    (3, 3, 1),
    (1076963, 4, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
-- +goose StatementEnd
