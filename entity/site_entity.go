package entity

import "time"

type Sites struct {
	Id, Domain, Title, Image, User_Id string
	Created_At, Updated_At            time.Time
	Links                             []Links
	Clicks                            int
}
