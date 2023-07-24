package models

import (
	"gorm.io/gorm"
	"time"
)

type Results struct {
	gorm.Model
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Results        string    `json:"results"`
	PartnerResults string    `json:"partnerResults"`
	Picture        string    `json:"picture"`
	PartnerPicture string    `json:"partnerPicture"`
	Date           time.Time `json:"time" `
}
