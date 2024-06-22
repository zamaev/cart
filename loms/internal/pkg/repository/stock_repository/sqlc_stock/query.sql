-- name: Reserve :exec
UPDATE stock
SET reserved = reserved + @count
WHERE sku = $1;

-- name: ReserveRemove :exec
UPDATE stock
SET
    total_count = total_count - @count,
    reserved = reserved - @count
WHERE sku = $1;

-- name: ReserveCancel :exec
UPDATE stock
SET reserved = reserved - @count
WHERE sku = $1;

-- name: GetStocksBySku :one
SELECT total_count - reserved AS count
FROM stock
WHERE sku = $1
LIMIT 1;
