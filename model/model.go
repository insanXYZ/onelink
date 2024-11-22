package model

import (
	"os"
	"radproject/entity"
)

// Convert entity to Response model
//
//	type UserResponse struct {
//		ID        string    `json:"id,omitempty"`
//		Name      string    `json:"name,omitempty"`
//		Email     string    `json:"email,omitempty"`
//		Password  string    `json:"password,omitempty"`
//		Image     string    `json:"image,omitempty"`
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		UpdatedAt time.Time `json:"updated_at,omitempty"`
//	}
func EntityToUserResponse(ent *entity.User) *UserResponse {
	return &UserResponse{
		Email: ent.Email,
		Name:  ent.Name,
		Image: "http://localhost" + os.Getenv("WEB_PORT") + "/storage/image/user/" + ent.Image,
		ID:    ent.ID,
	}
}
