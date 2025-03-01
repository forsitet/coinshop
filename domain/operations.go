package domain

import "time"

type Operations struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uint
	FromUser  string `json:"from_user"`
	ToUser    string `json:"to_user"`
	Amount    int    `json:"amount"`
}
