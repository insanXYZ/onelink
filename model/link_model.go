package model

import "time"

type LinkResponse struct {
	Id        string    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Href      string    `json:"href,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"Updated_at,omitempty"`
}

type CreateLinkRequest struct {
	Site_Id string `param:"id" validate:"required"`
	Title   string `form:"title" validate:"min=3,required"`
	Href    string `form:"href" validate:"required"`
}

type DeleteLinkRequest struct {
	Site_Id string `param:"site_id" validate:"required"`
	Link_Id string `param:"link_id" validate:"required"`
}
