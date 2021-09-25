package sql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-rel/rel"
)

// ErrorMapper function.
type ErrorMapper func(error) error

// IncrementFunc function.
type IncrementFunc func(SQL) int

// SQL base adapter.
type SQL struct {
	QueryBuilder     QueryBuilder
	InsertBuilder    InsertBuilder
	InsertAllBuilder InsertAllBuilder
	UpdateBuilder    UpdateBuilder
	DeleteBuilder    DeleteBuilder
	TableBuilder     TableBuilder
	IndexBuilder     IndexBuilder
	IncrementFunc    IncrementFunc
	ErrorMapper      ErrorMapper
	DB               *sql.DB
	Tx               *sql.Tx
	Savepoint        int
	Instrumenter     rel.Instrumenter
}

// Instrumentation set instrumenter for this adapter.
func (s *SQL) Instrumentation(instrumenter rel.Instrumenter) {
	s.Instrumenter = instrumenter
}

// DoExec using active database connection.
func (s SQL) DoExec(ctx context.Context, statement string, args []interface{}) (sql.Result, error) {
	var (
		err    error
		result sql.Result
		finish = s.Instrumenter.Observe(ctx, "adapter-exec", statement)
	)

	if s.Tx != nil {
		result, err = s.Tx.ExecContext(ctx, statement, args...)
	} else {
		result, err = s.DB.ExecContext(ctx, statement, args...)
	}

	finish(err)
	return result, err
}

// DoQuery using active database connection.
func (s SQL) DoQuery(ctx context.Context, statement string, args []interface{}) (*sql.Rows, error) {
	var (
		err  error
		rows *sql.Rows
	)

	finish := s.Instrumenter.Observe(ctx, "adapter-query", statement)
	if s.Tx != nil {
		rows, err = s.Tx.QueryContext(ctx, statement, args...)
	} else {
		rows, err = s.DB.QueryContext(ctx, statement, args...)
	}
	finish(err)

	return rows, err
}

// Begin begins a new transaction.
func (s SQL) Begin(ctx context.Context) (rel.Adapter, error) {
	var (
		tx        *sql.Tx
		savepoint int
		err       error
	)

	finish := s.Instrumenter.Observe(ctx, "adapter-begin", "begin transaction")

	if s.Tx != nil {
		tx = s.Tx
		savepoint = s.Savepoint + 1
		_, err = s.Tx.ExecContext(ctx, "SAVEPOINT s"+strconv.Itoa(savepoint)+";")
	} else {
		tx, err = s.DB.BeginTx(ctx, nil)
	}

	finish(err)

	return &SQL{
		QueryBuilder:     s.QueryBuilder,
		InsertBuilder:    s.InsertBuilder,
		InsertAllBuilder: s.InsertAllBuilder,
		UpdateBuilder:    s.UpdateBuilder,
		DeleteBuilder:    s.DeleteBuilder,
		TableBuilder:     s.TableBuilder,
		IndexBuilder:     s.IndexBuilder,
		IncrementFunc:    s.IncrementFunc,
		ErrorMapper:      s.ErrorMapper,
		Tx:               tx,
		Savepoint:        savepoint,
		Instrumenter:     s.Instrumenter,
	}, s.ErrorMapper(err)
}

// Commit commits current transaction.
func (s SQL) Commit(ctx context.Context) error {
	var err error

	finish := s.Instrumenter.Observe(ctx, "adapter-commit", "commit transaction")

	if s.Tx == nil {
		err = errors.New("unable to commit outside transaction")
	} else if s.Savepoint > 0 {
		_, err = s.Tx.ExecContext(ctx, "RELEASE SAVEPOINT s"+strconv.Itoa(s.Savepoint)+";")
	} else {
		err = s.Tx.Commit()
	}

	finish(err)

	return s.ErrorMapper(err)
}

// Rollback revert current transaction.
func (s SQL) Rollback(ctx context.Context) error {
	var err error

	finish := s.Instrumenter.Observe(ctx, "adapter-rollback", "rollback transaction")

	if s.Tx == nil {
		err = errors.New("unable to rollback outside transaction")
	} else if s.Savepoint > 0 {
		_, err = s.Tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT s"+strconv.Itoa(s.Savepoint)+";")
	} else {
		err = s.Tx.Rollback()
	}

	finish(err)

	return s.ErrorMapper(err)
}

// Ping database.
func (s SQL) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

// Close database connection.
//
// TODO: add closer to adapter interface
func (s SQL) Close() error {
	return s.DB.Close()
}

// Query performs query operation.
func (s SQL) Query(ctx context.Context, query rel.Query) (rel.Cursor, error) {
	var (
		statement, args = s.QueryBuilder.Build(query)
		rows, err       = s.DoQuery(ctx, statement, args)
	)

	return &Cursor{Rows: rows}, s.ErrorMapper(err)
}

// Exec performs exec operation.
func (s SQL) Exec(ctx context.Context, statement string, args []interface{}) (int64, int64, error) {
	var (
		res, err = s.DoExec(ctx, statement, args)
	)

	if err != nil {
		return 0, 0, s.ErrorMapper(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

// Aggregate record using given query.
func (s SQL) Aggregate(ctx context.Context, query rel.Query, mode string, field string) (int, error) {
	var (
		out             sql.NullInt64
		aggregateField  = "^" + mode + "(" + field + ") AS result"
		aggregateQuery  = query.Select(append([]string{aggregateField}, query.GroupQuery.Fields...)...)
		statement, args = s.QueryBuilder.Build(aggregateQuery)
		rows, err       = s.DoQuery(ctx, statement, args)
	)

	defer rows.Close()
	if err == nil && rows.Next() {
		rows.Scan(&out)
	}

	return int(out.Int64), s.ErrorMapper(err)
}

// Insert inserts a record to database and returns its id.
func (s SQL) Insert(ctx context.Context, query rel.Query, primaryField string, mutates map[string]rel.Mutate) (interface{}, error) {
	var (
		statement, args = s.InsertBuilder.Build(query.Table, primaryField, mutates)
		id, _, err      = s.Exec(ctx, statement, args)
	)

	return id, err
}

// InsertAll inserts multiple records to database and returns its ids.
func (s SQL) InsertAll(ctx context.Context, query rel.Query, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate) ([]interface{}, error) {
	var (
		statement, args = s.InsertAllBuilder.Build(query.Table, primaryField, fields, bulkMutates)
		id, _, err      = s.Exec(ctx, statement, args)
	)

	if err != nil {
		return nil, err
	}

	var (
		ids = make([]interface{}, len(bulkMutates))
		inc = 1
	)

	if s.IncrementFunc != nil {
		inc = s.IncrementFunc(s)
	}

	if inc < 0 {
		id = id + int64((len(bulkMutates)-1)*inc)
		inc *= -1
	}

	if primaryField != "" {
		counter := 0
		for i := range ids {
			if mut, ok := bulkMutates[i][primaryField]; ok {
				ids[i] = mut.Value
				id = toInt64(ids[i])
				counter = 1
			} else {
				ids[i] = id + int64(counter*inc)
				counter++
			}
		}
	}

	return ids, nil
}

// Update updates a record in database.
func (s SQL) Update(ctx context.Context, query rel.Query, primaryField string, mutates map[string]rel.Mutate) (int, error) {
	var (
		statement, args      = s.UpdateBuilder.Build(query.Table, primaryField, mutates, query.WhereQuery)
		_, updatedCount, err = s.Exec(ctx, statement, args)
	)

	return int(updatedCount), err
}

// Delete deletes all results that match the query.
func (s SQL) Delete(ctx context.Context, query rel.Query) (int, error) {
	var (
		statement, args      = s.DeleteBuilder.Build(query.Table, query.WhereQuery)
		_, deletedCount, err = s.Exec(ctx, statement, args)
	)

	return int(deletedCount), err
}

// SchemaApply performs migration to database.
func (s SQL) SchemaApply(ctx context.Context, migration rel.Migration) error {
	var (
		statement string
	)

	switch v := migration.(type) {
	case rel.Table:
		statement = s.TableBuilder.Build(v)
	case rel.Index:
		statement = s.IndexBuilder.Build(v)
	case rel.Raw:
		statement = string(v)
	}

	_, _, err := s.Exec(ctx, statement, nil)
	return err
}

// Apply performs migration to database.
//
// Deprecated: Use Schema Apply instead.
func (s SQL) Apply(ctx context.Context, migration rel.Migration) error {
	return s.SchemaApply(ctx, migration)
}
