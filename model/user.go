package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const SystemUserID = "system"

var systemUser = User{
	ID:   SystemUserID,
	Name: SystemUserID,
}

func SystemUser() User {
	return systemUser
}
