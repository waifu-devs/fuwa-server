package database

import (
	"context"
	"database/sql"
	"fmt"

	"go.opentelemetry.io/otel"
)

type TracedDB struct {
	db *sql.DB
}

func NewTracedDB(db *sql.DB) *TracedDB {
	return &TracedDB{db: db}
}

func (t *TracedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tracer := otel.Tracer("TracedDB")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("ExecContext: %s", query))
	defer span.End()

	return t.db.ExecContext(ctx, query, args...)
}

func (t *TracedDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	tracer := otel.Tracer("TracedDB")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("PrepareContext: %s", query))
	defer span.End()

	return t.db.PrepareContext(ctx, query)
}

func (t *TracedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tracer := otel.Tracer("TracedDB")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("QueryContext: %s", query))
	defer span.End()

	return t.db.QueryContext(ctx, query, args...)
}

func (t *TracedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tracer := otel.Tracer("TracedDB")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("QueryRowContext: %s", query))
	defer span.End()

	return t.db.QueryRowContext(ctx, query, args...)
}
