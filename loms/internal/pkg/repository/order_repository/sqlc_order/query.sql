-- name: Create :one
INSERT INTO orders
    (id, user_id, status)
VALUES
    (nextval('order_id_manual_seq') + @shard_id::int, $1, $2)
RETURNING id;

-- name: AddItem :exec
INSERT INTO order_items
    (order_id, sku, count)
VALUES
    ($1, $2, $3);

-- name: GetById :many
SELECT sqlc.embed(orders), sqlc.embed(order_items)
FROM orders
LEFT JOIN order_items ON orders.id = order_items.order_id
WHERE id = $1;

-- name: SetStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: GetAll :many
SELECT sqlc.embed(orders), sqlc.embed(order_items)
FROM orders
LEFT JOIN order_items ON orders.id = order_items.order_id;
