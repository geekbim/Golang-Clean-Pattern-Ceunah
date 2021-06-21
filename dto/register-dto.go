package dto

// RegisterDTO is used when client post from /register url
type RegisterDTO struct {
	Name     string `json:"name" form:"name" binding:"required" validate:"min:1"`
	Email    string `json:"email" form:"email" binding:"required" validate:"email"`
	Password string `jdon:"password" form:"password" binding:"required" validate:"min6"`
}
