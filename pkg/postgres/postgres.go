package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type IPostgres interface {
	Query(query string,
		args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExec(
		ctx context.Context,
		arg interface{},
		query string,
	) (res sql.Result, err error)
}

type Postgres struct {
	DB        *sqlx.DB
	cnnString *string
}

func New(connString string) (postgres *Postgres, err error) {
	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return postgres, err
	}

	postgres = &Postgres{
		DB:        db,
		cnnString: &connString,
	}
	return
}

func (p *Postgres) NamedExec(
	ctx context.Context,
	arg interface{},
	query string,
) (res sql.Result, err error) {
	stmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		return
	}
	res, err = stmt.ExecContext(ctx, arg)

	return
}

func (p *Postgres) Query(query string, args ...any) (*sql.Rows, error) {
	return p.DB.QueryContext(context.Background(), query, args...)
}

func (p *Postgres) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return p.DB.QueryContext(ctx, query, args...)
}
func (p *Postgres) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.DB.ExecContext(ctx, query, args...)
}
