package entity

import "time"

type Links struct {
	Id, Title, Href, Site_Id string
	CreatedAt, UpdatedAt     time.Time
	Clicks                   int
}
