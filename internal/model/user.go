package model

type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	Login    string `json:"login"`
	Password string `json:"password" gorm:"column:password_hash"`
}

type Users []*User

type ContextKey string

const ContextKeyUserID = ContextKey("user_id")
