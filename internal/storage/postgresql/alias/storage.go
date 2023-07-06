package alias

import (
	"database/sql"
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
		op    = "storage.postgresql.alias.SaveURL"
		query = "INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id"
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

func (s *Storage) DeleteAlias(alias string) error {
	const (
		op    = "storage.postgresql.alias.DeleteURL"
		query = "DELETE FROM url WHERE alias = $1"
	)

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
