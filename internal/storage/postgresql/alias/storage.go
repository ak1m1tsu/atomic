package alias

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/romankravchuk/atomic/internal/data"
	"github.com/romankravchuk/atomic/internal/storage"
	"github.com/romankravchuk/atomic/internal/storage/postgresql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	const op = "storage.postgresql.alias.New"

	if db == nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidArgument)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveAlias(alias *data.Alias) error {
	const (
		op    = "storage.postgresql.alias.SaveAlias"
		query = "INSERT INTO alias (url, name) VALUES ($1, $2) RETURNING id"
	)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(alias.URL, alias.Name).Scan(&alias.ID)
	if err != nil {
		if psqlErr, ok := err.(*pq.Error); ok && psqlErr.Code == postgresql.UniqueViolation {
			return fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteAlias(name string) error {
	const (
		op    = "storage.postgresql.alias.DeleteAlias"
		query = "DELETE FROM alias WHERE name = $1"
	)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if count != 1 {
		return fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
	}

	return nil
}

func (s *Storage) GetAlias(name string) (string, error) {
	const (
		op    = "storage.postgresql.alias.GetAlias"
		query = "SELECT url FROM alias WHERE name = $1"
	)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string

	err = stmt.QueryRow(name).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
