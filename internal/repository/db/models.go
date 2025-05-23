// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package dbCtx

import (
	"database/sql"
)

type Message struct {
	ID        int32        `db:"id" json:"id"`
	SenderID  int32        `db:"sender_id" json:"senderId"`
	Content   string       `db:"content" json:"content"`
	CreatedAt sql.NullTime `db:"created_at" json:"createdAt"`
}

type SchemaMigration struct {
	Version string `db:"version" json:"version"`
}

type User struct {
	ID           int32        `db:"id" json:"id"`
	Username     string       `db:"username" json:"username"`
	Email        string       `db:"email" json:"email"`
	FullName     string       `db:"full_name" json:"fullName"`
	PasswordHash string       `db:"password_hash" json:"passwordHash"`
	CreatedAt    sql.NullTime `db:"created_at" json:"createdAt"`
	UpdatedAt    sql.NullTime `db:"updated_at" json:"updatedAt"`
	DeletedAt    sql.NullTime `db:"deleted_at" json:"deletedAt"`
}
