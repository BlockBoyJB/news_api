package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

const (
	defaultMaxPoolSize = 10
)

type Postgres interface {
	SelectRows(view reform.View, tail string, args ...interface{}) (*sql.Rows, error)
	Begin() (*reform.TX, error)
	Close()
}

type postgres struct {
	*reform.DB
	conn        *sql.DB
	maxPoolSize int
}

func NewPG(url string, opts ...Option) (Postgres, error) {
	pg := &postgres{
		maxPoolSize: defaultMaxPoolSize,
	}
	for _, option := range opts {
		option(pg)
	}
	conn, err := sql.Open("postgres", url+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(pg.maxPoolSize)
	pg.conn = conn

	pg.DB = reform.NewDB(conn, postgresql.Dialect, nil)
	return pg, nil
}

func (p *postgres) Close() {
	_ = p.conn.Close()
}
