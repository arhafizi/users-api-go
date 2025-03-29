package dto

import "time"

type UserResponse struct {
	ID        int32      `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"fullName"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
