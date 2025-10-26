package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	// ErrNoRows is returned when no rows are found in the database.
	ErrNoRows           = pgx.ErrNoRows
	ErrEmptyFilterAttrs = errors.New("no filter attributes provided")
	ErrNoUpdateAttrs    = errors.New("no update attributes provided")
)
