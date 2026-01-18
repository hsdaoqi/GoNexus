package vo

import "time"

type UpdateUser struct {
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Location  string    `json:"Location"`
	Signature string    `json:"signature"`
}
