// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package dbCtx

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, email, full_name, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email, full_name, password_hash, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	FullName     string `db:"full_name" json:"fullName"`
	PasswordHash string `db:"password_hash" json:"passwordHash"`
}

// queries.sql
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, email, full_name, password_hash, created_at, updated_at, deleted_at FROM users
WHERE email = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, full_name, password_hash, created_at, updated_at, deleted_at FROM users
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, full_name, password_hash, created_at, updated_at, deleted_at FROM users
WHERE username = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, full_name, password_hash, created_at, updated_at, deleted_at FROM users
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.FullName,
			&i.PasswordHash,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteUser = `-- name: SoftDeleteUser :execrows
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) SoftDeleteUser(ctx context.Context, id int32) (int64, error) {
	result, err := q.db.ExecContext(ctx, softDeleteUser, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateUserFull = `-- name: UpdateUserFull :one
UPDATE users
SET username = $2, email = $3, full_name = $4, password_hash = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, username, email, full_name, password_hash, created_at, updated_at, deleted_at
`

type UpdateUserFullParams struct {
	ID           int32  `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	FullName     string `db:"full_name" json:"fullName"`
	PasswordHash string `db:"password_hash" json:"passwordHash"`
}

func (q *Queries) UpdateUserFull(ctx context.Context, arg UpdateUserFullParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserFull,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUserPartial = `-- name: UpdateUserPartial :one
UPDATE users
SET
    username = COALESCE($2, username),
    email = COALESCE($3, email),
    full_name = COALESCE($4, full_name),
    password_hash = COALESCE($5, password_hash),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, username, email, full_name, password_hash, created_at, updated_at, deleted_at
`

type UpdateUserPartialParams struct {
	ID           int32  `db:"id" json:"id"`
	Username     *string `db:"username" json:"username"`
	Email        *string `db:"email" json:"email"`
	FullName     *string `db:"full_name" json:"fullName"`
	PasswordHash *string `db:"password_hash" json:"passwordHash"`
}

func (q *Queries) UpdateUserPartial(ctx context.Context, arg UpdateUserPartialParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPartial,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FullName,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
