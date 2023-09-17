package model

type CheckinReq struct {
	Uid  string `json:"uid" binding:"required"`
	Time string `json:"time" binding:"required"` // 2023-09-17 07:00:12
}

type User struct {
	Uid      string `json:"uid" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

type Checkin struct {
	Uid  string `json:"uid"`
	Time string `json:"time" binding:"required"` // 07:00:12
}
