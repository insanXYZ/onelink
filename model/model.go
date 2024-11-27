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
		Image: os.Getenv("APP_URL") + imageDir + "user/" + ent.Image,
		ID:    ent.ID,
	}
}

func EntityToSiteResponse(ent *entity.Sites) *SiteResponse {
	res := &SiteResponse{
		Id:         ent.Id,
		Domain:     ent.Domain,
		Title:      ent.Title,
		Image:      os.Getenv("APP_URL") + imageDir + "site/" + ent.Image,
		Created_At: ent.Created_At,
		Updated_At: ent.Updated_At,
		Clicks:     ent.Clicks,
	}

	if len(ent.Links) > 0 {
		for _, v := range ent.Links {
			if v.Id != "" {
				res.Links = append(res.Links, *EntitytoLinkResponse(&v))
			}
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
		Clicks:    ent.Clicks,
	}
}
