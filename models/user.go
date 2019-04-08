package models

// Информация о пользователе.
type User struct {
	About string `json:"about,omitempty"`
	Email string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

type Users []User