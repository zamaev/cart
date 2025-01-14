// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlc_outbox

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const create = `-- name: Create :one
INSERT INTO outbox
    (topic, event, headers)
VALUES
    ($1, $2, $3)
RETURNING id
`

type CreateParams struct {
	Topic   string
	Event   []byte
	Headers []byte
}

func (q *Queries) Create(ctx context.Context, arg CreateParams) (int64, error) {
	row := q.db.QueryRow(ctx, create, arg.Topic, arg.Event, arg.Headers)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getWaitList = `-- name: GetWaitList :many
SELECT id, topic, event, headers, created_at
FROM outbox
WHERE completed_at IS NULL
`

type GetWaitListRow struct {
	ID        int64
	Topic     string
	Event     []byte
	Headers   []byte
	CreatedAt pgtype.Timestamp
}

func (q *Queries) GetWaitList(ctx context.Context) ([]GetWaitListRow, error) {
	rows, err := q.db.Query(ctx, getWaitList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWaitListRow
	for rows.Next() {
		var i GetWaitListRow
		if err := rows.Scan(
			&i.ID,
			&i.Topic,
			&i.Event,
			&i.Headers,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setComplete = `-- name: SetComplete :exec
UPDATE outbox
SET completed_at = now()
WHERE id = $1
`

func (q *Queries) SetComplete(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, setComplete, id)
	return err
}
