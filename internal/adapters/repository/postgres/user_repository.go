package postgres

import (
	"context"
	"crud/internal/domain/entities"
	"crud/internal/services/user"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, attrs entities.UserAttrs, ent *entities.User) error {
	const insert = `
		INSERT INTO users (id, username, email, hashed_password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, hashed_password`

	if ctx.Err() != nil {
		return ctx.Err()
	}

	row := r.pool.QueryRow(ctx, insert, attrs.ID, attrs.Username, attrs.Email, attrs.HashedPassword)

	err := row.Scan(&ent.ID, &ent.Username, &ent.Email, &ent.HashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_username_key" {
					return user.ErrUsernameTaken
				}
				if pgErr.ConstraintName == "users_email_key" {
					return user.ErrEmailTaken
				}
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindOne(ctx context.Context, filterAttrs entities.UserFilterAttrs, ent *entities.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	clauses := []string{}
	args := []any{}

	if v, ok := filterAttrs.ID.Get(); ok {
		args = append(args, v)
		clauses = append(clauses, fmt.Sprintf("id = $%d", len(args)))
	}
	if v, ok := filterAttrs.Email.Get(); ok {
		args = append(args, v)
		clauses = append(clauses, fmt.Sprintf("email = $%d", len(args)))
	}
	if v, ok := filterAttrs.Username.Get(); ok {
		args = append(args, v)
		clauses = append(clauses, fmt.Sprintf("username = $%d", len(args)))
	}
	if len(clauses) == 0 {
		return errors.New("empty filter")
	}

	insert := fmt.Sprintf(`SELECT id, username, email, hashed_password FROM users WHERE %s LIMIT 1`, strings.Join(clauses, "AND"))

	row := r.pool.QueryRow(ctx, insert, args...)

	if err := row.Scan(&ent.ID, &ent.Username, &ent.Email, &ent.HashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepository) Find(ctx context.Context, filterAttrs entities.UserFilterAttrs, ent *entities.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	clauses := []string{}
	args := []any{}

	if v, ok := filterAttrs.Email.Get(); ok {
		args = append(args, v)
		clauses = append(clauses, fmt.Sprintf("email = $%d", len(args)))
	}
	if v, ok := filterAttrs.Username.Get(); ok {
		args = append(args, v)
		clauses = append(clauses, fmt.Sprintf("username = $%d", len(args)))
	}
	if len(clauses) == 0 {
		return errors.New("empty filter")
	}

	insert := fmt.Sprintf(`SELECT id, username, email, hashed_password FROM users WHERE %s`, strings.Join(clauses, "AND"))

	row := r.pool.QueryRow(ctx, insert, args...)

	if err := row.Scan(&ent.ID, &ent.Username, &ent.Email, &ent.HashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}
	return nil
}