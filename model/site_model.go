package model

import "mime/multipart"

type CreateSiteRequest struct {
	Title string `form:"title" validate:"min=3,required"`
	Image *multipart.FileHeader
}
