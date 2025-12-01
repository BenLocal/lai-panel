package database

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/ngrok/sqlmw"
)

var DB *sqlx.DB

func InitDB(dbPath string) error {
	var err error
	db, err := sqlx.Connect("sqlite3-mw", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	db.MustExec("PRAGMA foreign_keys = ON")
	db.MustExec("PRAGMA journal_mode = WAL")

	err = runMigrations(db)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrations up to date")
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	DB = db
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sqlx.DB {
	return DB
}

type sqlInterceptor struct {
	sqlmw.NullInterceptor
}

func (in *sqlInterceptor) StmtQueryContext(ctx context.Context, conn driver.StmtQueryContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	startedAt := time.Now()
	rows, err := conn.QueryContext(ctx, args)
	log.Printf("[SQL] stmt_query err=%v duration=%s\n%s", err, time.Since(startedAt), formatQuery(query, args))
	return ctx, rows, err
}

func (in *sqlInterceptor) StmtExecContext(ctx context.Context, conn driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	startedAt := time.Now()
	result, err := conn.ExecContext(ctx, args)
	log.Printf("[SQL] stmt_exec err=%v duration=%s\n%s", err, time.Since(startedAt), formatQuery(query, args))
	return result, err
}

func (in *sqlInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	startedAt := time.Now()
	result, err := conn.ExecContext(ctx, query, args)
	log.Printf("[SQL] exec err=%v duration=%s\n%s", err, time.Since(startedAt), formatQuery(query, args))
	return result, err
}

func (in *sqlInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	startedAt := time.Now()
	rows, err := conn.QueryContext(ctx, query, args)
	log.Printf("[SQL] query err=%v duration=%s\n%s", err, time.Since(startedAt), formatQuery(query, args))
	return ctx, rows, err
}

func (in *sqlInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (context.Context, driver.Stmt, error) {
	startedAt := time.Now()
	stmt, err := conn.PrepareContext(ctx, query)
	log.Printf("[SQL] prepare err=%v duration=%s\n%s", err, time.Since(startedAt), query)
	return ctx, stmt, err
}

func (in *sqlInterceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (context.Context, driver.Tx, error) {
	startedAt := time.Now()
	tx, err := conn.BeginTx(ctx, txOpts)
	log.Printf("[SQL] begin_tx options=%+v err=%v duration=%s", txOpts, err, time.Since(startedAt))
	return ctx, tx, err
}

func (in *sqlInterceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) error {
	startedAt := time.Now()
	err := rows.Next(dest)
	log.Printf("[SQL] rows_next err=%v duration=%s", err, time.Since(startedAt))
	return err
}

func (in *sqlInterceptor) RowsClose(ctx context.Context, rows driver.Rows) error {
	startedAt := time.Now()
	err := rows.Close()
	log.Printf("[SQL] rows_close err=%v duration=%s", err, time.Since(startedAt))
	return err
}

func (in *sqlInterceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	startedAt := time.Now()
	err := tx.Commit()
	log.Printf("[SQL] tx_commit err=%v duration=%s", err, time.Since(startedAt))
	return err
}

func (in *sqlInterceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	startedAt := time.Now()
	err := tx.Rollback()
	log.Printf("[SQL] tx_rollback err=%v duration=%s", err, time.Since(startedAt))
	return err
}

func (in *sqlInterceptor) StmtClose(ctx context.Context, stmt driver.Stmt) error {
	startedAt := time.Now()
	err := stmt.Close()
	log.Printf("[SQL] stmt_close err=%v duration=%s", err, time.Since(startedAt))
	return err
}

func init() {
	sql.Register("sqlite3-mw", sqlmw.Driver(&sqlite3.SQLiteDriver{}, &sqlInterceptor{}))
}

func formatQuery(query string, args []driver.NamedValue) string {
	if len(args) == 0 {
		return query
	}

	replacer := strings.NewReplacer("?", "%s")
	placeholderCount := strings.Count(query, "?")

	values := make([]interface{}, 0, len(args))
	for _, arg := range args {
		values = append(values, quoteValue(arg.Value))
	}

	if placeholderCount != len(values) {
		var buf bytes.Buffer
		buf.WriteString(query)
		buf.WriteString("\n-- args: ")
		for i, value := range values {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf("%v", value))
		}
		return buf.String()
	}

	return fmt.Sprintf(replacer.Replace(query), values...)
}

func quoteValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format(time.RFC3339Nano))
	case []byte:
		if len(v) == 0 {
			return "x''"
		}
		return fmt.Sprintf("x'%X'", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
