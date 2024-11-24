package model

import (
	"os"
	"radproject/entity"
)

const imageDir = "/storage/image/"

func EntityToUserResponse(ent *entity.User) *UserResponse {
	return &UserResponse{
		Email: ent.Email,
		Name:  ent.Name,
		Image: "http://localhost" + os.Getenv("WEB_PORT") + imageDir + "user/" + ent.Image,
		ID:    ent.ID,
	}
}

func EntityToSiteResponse(ent *entity.Sites) *SiteResponse {
	return &SiteResponse{
		Id:     ent.Id,
		Domain: ent.Domain,
		Title:  ent.Title,
		Image:  "http://localhost" + os.Getenv("WEB_PORT") + imageDir + "site/" + ent.Image,
	}
}
