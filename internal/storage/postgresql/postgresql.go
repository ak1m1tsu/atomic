package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const UniqueViolation = "23505"

func New(dsn string) (*sql.DB, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
