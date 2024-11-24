package model

import (
	"mime/multipart"
	"time"
)

type SiteResponse struct {
	Id         string    `json:"id,omitempty"`
	Domain     string    `json:"domain,omitempty"`
	Title      string    `json:"title,omitempty"`
	Image      string    `json:"image,omitempty"`
	User_Id    string    `json:"user_id,omitempty"`
	Created_At time.Time `json:"created_at,omitempty"`
	Updated_At time.Time `json:"updated_at,omitempty"`
}

type CreateSiteRequest struct {
	Title  string `form:"title" validate:"min=3,required"`
	Domain string `form:"domain" validate:"min=3,required"`
	Image  *multipart.FileHeader
}

type DeleteSiteRequest struct {
	Id string `form:"id" validate:"required"`
}
