package models

type Telegram struct {
	Model
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}
