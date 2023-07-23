package model

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Users []*User
