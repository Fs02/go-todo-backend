// Package postgres wraps postgres (pq) driver as an adapter for REL.
//
// Usage:
//	// open postgres connection.
//	adapter, err := postgres.Open("postgres://postgres@localhost/rel_test?sslmode=disable")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize REL's repo.
//	repo := rel.New(adapter)
package postgres

import (
	"context"
	db "database/sql"
	"time"

	"github.com/go-rel/rel"
	"github.com/go-rel/sql"
	"github.com/go-rel/sql/builder"
)

// Postgres adapter.
type Postgres struct {
	sql.SQL
}

// New postgres adapter using existing connection.
func New(database *db.DB) rel.Adapter {
	var (
		bufferFactory     = builder.BufferFactory{ArgumentPlaceholder: "$", ArgumentOrdinal: true, BoolTrueValue: "true", BoolFalseValue: "false", Quoter: Quote{}, ValueConverter: ValueConvert{}}
		filterBuilder     = builder.Filter{}
		queryBuilder      = builder.Query{BufferFactory: bufferFactory, Filter: filterBuilder}
		OnConflictBuilder = builder.OnConflict{Statement: "ON CONFLICT", IgnoreStatement: "DO NOTHING", UpdateStatement: "DO UPDATE SET", TableQualifier: "excluded", SupportKey: true}
		InsertBuilder     = builder.Insert{BufferFactory: bufferFactory, ReturningPrimaryValue: true, InsertDefaultValues: true, OnConflict: OnConflictBuilder}
		insertAllBuilder  = builder.InsertAll{BufferFactory: bufferFactory, ReturningPrimaryValue: true, OnConflict: OnConflictBuilder}
		updateBuilder     = builder.Update{BufferFactory: bufferFactory, Query: queryBuilder, Filter: filterBuilder}
		deleteBuilder     = builder.Delete{BufferFactory: bufferFactory, Query: queryBuilder, Filter: filterBuilder}
		ddlBufferFactory  = builder.BufferFactory{InlineValues: true, BoolTrueValue: "true", BoolFalseValue: "false", Quoter: Quote{}, ValueConverter: ValueConvert{}}
		ddlQueryBuilder   = builder.Query{BufferFactory: ddlBufferFactory, Filter: filterBuilder}
		tableBuilder      = builder.Table{BufferFactory: ddlBufferFactory, ColumnMapper: columnMapper}
		indexBuilder      = builder.Index{BufferFactory: ddlBufferFactory, Query: ddlQueryBuilder, Filter: filterBuilder, SupportFilter: true}
	)

	return &Postgres{
		SQL: sql.SQL{
			QueryBuilder:     queryBuilder,
			InsertBuilder:    InsertBuilder,
			InsertAllBuilder: insertAllBuilder,
			UpdateBuilder:    updateBuilder,
			DeleteBuilder:    deleteBuilder,
			TableBuilder:     tableBuilder,
			IndexBuilder:     indexBuilder,
			ErrorMapper:      errorMapper,
			DB:               database,
		},
	}
}

// Open postgres connection using dsn.
func Open(dsn string) (rel.Adapter, error) {
	var database, err = db.Open("postgres", dsn)
	return New(database), err
}

// MustOpen postgres connection using dsn.
func MustOpen(dsn string) rel.Adapter {
	var database, err = db.Open("postgres", dsn)
	check(err)
	return New(database)
}

// Insert inserts a record to database and returns its id.
func (p Postgres) Insert(ctx context.Context, query rel.Query, primaryField string, mutates map[string]rel.Mutate, onConflict rel.OnConflict) (interface{}, error) {
	var (
		id              int64
		statement, args = p.InsertBuilder.Build(query.Table, primaryField, mutates, onConflict)
		rows, err       = p.DoQuery(ctx, statement, args)
	)

	if err == nil && rows.Next() {
		defer rows.Close()
		rows.Scan(&id)
	}

	return id, p.ErrorMapper(err)
}

// InsertAll inserts multiple records to database and returns its ids.
func (p Postgres) InsertAll(ctx context.Context, query rel.Query, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate, onConflict rel.OnConflict) ([]interface{}, error) {
	var (
		ids             []interface{}
		statement, args = p.InsertAllBuilder.Build(query.Table, primaryField, fields, bulkMutates, onConflict)
		rows, err       = p.DoQuery(ctx, statement, args)
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int64
			rows.Scan(&id)
			ids = append(ids, id)
		}
	}

	return ids, p.ErrorMapper(err)
}

// Begin begins a new transaction.
func (p Postgres) Begin(ctx context.Context) (rel.Adapter, error) {
	var (
		txSql, err = p.SQL.Begin(ctx)
	)

	return &Postgres{SQL: *txSql.(*sql.SQL)}, err
}

func errorMapper(err error) error {
	if err == nil {
		return nil
	}

	var (
		msg            = err.Error()
		constraintType = sql.ExtractString(msg, "violates ", " constraint")
	)

	switch constraintType {
	case "unique":
		return rel.ConstraintError{
			Key:  sql.ExtractString(err.Error(), "constraint \"", "\""),
			Type: rel.UniqueConstraint,
			Err:  err,
		}
	case "foreign key":
		return rel.ConstraintError{
			Key:  sql.ExtractString(err.Error(), "constraint \"", "\""),
			Type: rel.ForeignKeyConstraint,
			Err:  err,
		}
	case "check":
		return rel.ConstraintError{
			Key:  sql.ExtractString(err.Error(), "constraint \"", "\""),
			Type: rel.CheckConstraint,
			Err:  err,
		}
	default:
		return err
	}
}

func columnMapper(column *rel.Column) (string, int, int) {
	var (
		typ  string
		m, n int
	)

	// postgres specific
	column.Unsigned = false
	if column.Default == "" {
		column.Default = nil
	}

	switch column.Type {
	case rel.ID:
		typ = "SERIAL NOT NULL"
	case rel.BigID:
		typ = "BIGSERIAL NOT NULL"
	case rel.DateTime:
		typ = "TIMESTAMPTZ"
		if t, ok := column.Default.(time.Time); ok {
			column.Default = FormatTime(t)
		}
	case rel.Int, rel.BigInt, rel.Text:
		column.Limit = 0
		typ, m, n = sql.ColumnMapper(column)
	case rel.JSON:
		typ = "JSONB"
	default:
		typ, m, n = sql.ColumnMapper(column)
	}

	return typ, m, n
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
