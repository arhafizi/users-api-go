package dto

type ListUsersParams struct {
	Limit  int32 `binding:"required,min=1,max=100"`
	Offset int32 `binding:"min=0"`
}
