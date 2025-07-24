package models

import "gorm.io/gorm"

type ForgotPassword struct {
	gorm.Model
	Token  string `json:"token" gorm:"unique"`
	UserId uint   `json:"userId"`
}

type ForgotPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
