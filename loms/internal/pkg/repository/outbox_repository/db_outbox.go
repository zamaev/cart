package outboxrepository

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/model"
	"route256/loms/internal/pkg/repository"
	"route256/loms/internal/pkg/repository/outbox_repository/sqlc_outbox"
	"route256/loms/pkg/tracing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DbOutboxRepository struct {
	db      DB
	queries *sqlc_outbox.Queries
}

func NewDbOutboxRepository(db DB) *DbOutboxRepository {
	return &DbOutboxRepository{
		db:      db,
		queries: sqlc_outbox.New(db),
	}
}

func (r *DbOutboxRepository) Create(ctx context.Context, topic string, event []byte, headers []byte) (_ int64, err error) {
	ctx, span := tracing.Start(ctx, "DbOutboxRepository.Create")
	defer tracing.EndWithCheckError(span, &err)

	queries := r.queries
	if tx, ok := ctx.Value(repository.CtxTxKey{}).(pgx.Tx); ok {
		queries = queries.WithTx(tx)
	}

	id, err := queries.Create(ctx, sqlc_outbox.CreateParams{
		Topic:   topic,
		Event:   event,
		Headers: headers,
	})
	if err != nil {
		return 0, fmt.Errorf("r.queries.Create: %w", err)
	}
	return id, nil
}

func (r *DbOutboxRepository) GetWaitList(ctx context.Context) (_ []model.OutboxItem, err error) {
	ctx, span := tracing.Start(ctx, "DbOutboxRepository.GetWaitList")
	defer tracing.EndWithCheckError(span, &err)

	rows, err := r.queries.GetWaitList(ctx)
	if err != nil {
		return nil, fmt.Errorf("r.queries.GetWaitList: %w", err)
	}

	items := make([]model.OutboxItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.OutboxItem{
			Id:        row.ID,
			Topic:     row.Topic,
			Event:     row.Event,
			Headers:   row.Headers,
			CreatedAt: row.CreatedAt.Time,
		})
	}
	return items, nil
}

func (r *DbOutboxRepository) SetComplete(ctx context.Context, id int64) (err error) {
	ctx, span := tracing.Start(ctx, "DbOutboxRepository.SetComplete")
	defer tracing.EndWithCheckError(span, &err)

	if err = r.queries.SetComplete(ctx, id); err != nil {
		return fmt.Errorf("r.queries.SetComplete: %w", err)
	}
	return nil
}
