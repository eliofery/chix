package model

const SessionTableName = "sessions"

type Session struct {
	ID     int64  `json:"id"`
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}
