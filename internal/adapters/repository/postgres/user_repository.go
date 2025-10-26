package postgres

import (
	"context"
	"crud/internal/domain/entities"
	"crud/internal/services/user"
	"errors"
	"fmt"
	"strings"

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
		return ErrEmptyFilterAttrs
	}

	query := fmt.Sprintf(`SELECT id, username, email, hashed_password FROM users WHERE %s LIMIT 1`, strings.Join(clauses, " AND "))

	row := r.pool.QueryRow(ctx, query, args...)

	if err := row.Scan(&ent.ID, &ent.Username, &ent.Email, &ent.HashedPassword); err != nil {
		if errors.Is(err, ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepository) Find(ctx context.Context, filterAttrs entities.UserFilterAttrs) ([]entities.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
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
		return nil, ErrEmptyFilterAttrs
	}

	query := fmt.Sprintf(`SELECT id, username, email, hashed_password FROM users WHERE %s`, strings.Join(clauses, " AND "))

	var ents []entities.User
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ent entities.User
		if err = rows.Scan(&ent.ID, &ent.Username, &ent.Email, &ent.HashedPassword); err != nil {
			return nil, err
		}
		ents = append(ents, ent)
	}
	return ents, nil
}

func (r *UserRepository) Update(ctx context.Context, attrs entities.UserUpdateAttrs, filterAttrs entities.UserFilterAttrs, ent *entities.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	filterClauses := []string{}
	args := []any{}

	if v, ok := filterAttrs.ID.Get(); ok {
		args = append(args, v)
		filterClauses = append(filterClauses, fmt.Sprintf("id = $%d", len(args)))
	}
	if v, ok := filterAttrs.Email.Get(); ok {
		args = append(args, v)
		filterClauses = append(filterClauses, fmt.Sprintf("email = $%d", len(args)))
	}
	if v, ok := filterAttrs.Username.Get(); ok {
		args = append(args, v)
		filterClauses = append(filterClauses, fmt.Sprintf("username = $%d", len(args)))
	}

	if len(filterClauses) == 0 {
		return ErrEmptyFilterAttrs
	}

	attrsClauses := []string{}

	if v, ok := attrs.Email.Get(); ok {
		args = append(args, v)
		attrsClauses = append(attrsClauses, fmt.Sprintf("email = $%d", len(args)))
	}
	if v, ok := attrs.Username.Get(); ok {
		args = append(args, v)
		attrsClauses = append(attrsClauses, fmt.Sprintf("username = $%d", len(args)))
	}
	if v, ok := attrs.HashedPassword.Get(); ok {
		args = append(args, v)
		attrsClauses = append(attrsClauses, fmt.Sprintf("hashed_password = $%d", len(args)))
	}

	if len(attrsClauses) == 0 {
		return ErrNoUpdateAttrs
	}

	query := fmt.Sprintf(`UPDATE users SET %s WHERE %s RETURNING id, username, email, hashed_password`,
		strings.Join(attrsClauses, ", "),
		strings.Join(filterClauses, " AND "))

	row := r.pool.QueryRow(ctx, query, args...)

	err := row.Scan(&ent.ID, &ent.Username, &ent.Email, &ent.HashedPassword)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			return user.ErrUserNotFound
		}
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
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, filterAttrs entities.UserFilterAttrs) error {
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
		return ErrEmptyFilterAttrs
	}

	query := fmt.Sprintf(`DELETE FROM users WHERE %s`, strings.Join(clauses, " AND "))
	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}
	return nil
}
