package entity

import "time"

type Sites struct {
	Id, Title, Image, User_Id string
	Created_At, Updated_At    time.Time
}
