package model

type Winner struct {
	User   User `json:"user"`
	Points int  `json:"points"`
}
