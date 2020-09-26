package repository

// Repository repository interface
type Repository interface {
	PostgresRepository
}

// PostgresRepository is a postgres repository interface
type PostgresRepository interface {
	PostgresRepositoryConn
	Close()
}

// PostgresRepositoryConn is a postgres connections interface
type PostgresRepositoryConn interface {
	GetConn() (*PostgresConn, error)
	CloseConn(c *PostgresConn)
}

// New create a new repository
func New(url string) Repository {
	r := NewPostgres(url)
	return r
}
