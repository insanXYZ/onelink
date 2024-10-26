package model

type LoginRequest struct {
	Email    string `validate:"email,required" form:"email"`
	Password string `validate:"min=6,required" form:"password"`
}

type RegisterRequest struct {
	Name     string `validate:"min=3,required" form:"name"`
	Email    string `validate:"email,required" form:"email"`
	Password string `validate:"min=6,required" form:"password"`
}
