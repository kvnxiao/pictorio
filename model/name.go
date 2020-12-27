package model

type NameRequest struct {
	Name string `json:"name"`
}

type NameResponse struct {
	Name        string `json:"name"`
	IsGenerated bool   `json:"isGenerated"`
}
