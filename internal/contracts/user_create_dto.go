package dto

type CreateUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required,max=100"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
