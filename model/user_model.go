package model

import (
	"mime/multipart"
	"time"
)

type UserResponse struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`
	Email        string         `json:"email,omitempty"`
	Password     string         `json:"password,omitempty"`
	Image        string         `json:"image,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
	SiteResponse []SiteResponse `json:"sites,omitempty"`
	LinkResponse []LinkResponse `json:"links,omitempty"`
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

type UpdateUserRequest struct {
	Name  string `validate:"omitempty,min=3" form:"name"`
	Email string `validate:"omitempty,email" form:"email"`
	Image *multipart.FileHeader
}

type DashboardRequest struct {
	From string `query:"from"`
	To   string `query:"to"`
}
