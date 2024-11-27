package model

import (
	"mime/multipart"
	"time"
)

type SiteResponse struct {
	Id         string         `json:"id,omitempty"`
	Domain     string         `json:"domain,omitempty"`
	Title      string         `json:"title,omitempty"`
	Image      string         `json:"image,omitempty"`
	User_Id    string         `json:"user_id,omitempty"`
	Links      []LinkResponse `json:"links,omitempty"`
	Created_At time.Time      `json:"created_at,omitempty"`
	Updated_At time.Time      `json:"updated_at,omitempty"`
	Clicks     int            `json:"clicks,omitempty"`
}

type ViewPublishSite struct {
	DomainSite string `param:"domain_site" validate:"required"`
}

type CreateSiteRequest struct {
	Title  string `form:"title" validate:"min=3,required"`
	Domain string `form:"domain" validate:"min=3,required"`
	Image  *multipart.FileHeader
}

type DeleteSiteRequest struct {
	Id string `param:"id" validate:"required"`
}

type UpdateSiteRequest struct {
	Id     string `param:"id" validate:"required"`
	Domain string `form:"domain" validate:"omitempty,min=3"`
	Title  string `form:"title" validate:"omitempty,min=3"`
	Image  *multipart.FileHeader
}
