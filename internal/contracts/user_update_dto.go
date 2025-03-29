package dto

type UpdateUserFullReq struct {
	ID       int32
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserPartialReq struct {
	ID       int32
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	FullName *string `json:"fullName,omitempty" binding:"omitempty,max=100"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
}
