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
	res := &SiteResponse{
		Id:         ent.Id,
		Domain:     ent.Domain,
		Title:      ent.Title,
		Image:      "http://localhost" + os.Getenv("WEB_PORT") + imageDir + "site/" + ent.Image,
		Created_At: ent.Created_At,
		Updated_At: ent.Updated_At,
	}

	if len(ent.Links) > 0 {
		for _, v := range ent.Links {
			res.Links = append(res.Links, *EntitytoLinkResponse(&v))
		}
	}

	return res
}

func EntitytoLinkResponse(ent *entity.Links) *LinkResponse {
	return &LinkResponse{
		Id:        ent.Id,
		Title:     ent.Title,
		Href:      ent.Href,
		UpdatedAt: ent.UpdatedAt,
		CreatedAt: ent.CreatedAt,
	}
}
