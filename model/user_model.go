package model

import "time"

type UserResponse struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Image     string    `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type LoginRequest struct {
	Email    string `validate:"email,required" form:"email"`
	Password string `validate:"min=6,required" form:"password"`
}

type RegisterRequest struct {
	Name     string `validate:"min=3,required" form:"name"`
	Email    string `validate:"email,required" form:"email"`
	Password string `validate:"min=6,required" form:"password"`
}
