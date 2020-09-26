package repository

import (
	"context"

	// postgres driver
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// Postgres postgres repo
type Postgres struct {
	pool *pgxpool.Pool
}

// PostgresConn postgres conn repo
type PostgresConn struct {
	Conn *pgxpool.Conn
}

// NewPostgres return a new postgres repository
func NewPostgres(url string) *Postgres {
	r := Postgres{}
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("repository.NewPostgresRepository(): url=%+v pool=%+v, error=%w", url, pool, err)
	}

	r.pool = pool

	return &r
}

// GetConn  get a connection from repository
func (r *Postgres) GetConn() (*PostgresConn, error) {
	p := PostgresConn{}
	var err error

	p.Conn, err = r.pool.Acquire(context.Background())

	return &p, err
}

// CloseConn close the connection acquired on postgres repository
func (r *Postgres) CloseConn(c *PostgresConn) {
	c.Conn.Release()
}

// Close close the postgres repository
func (r *Postgres) Close() {
	r.pool.Close()
}
