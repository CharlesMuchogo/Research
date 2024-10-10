package models

import "github.com/google/uuid"

type Results struct {
	Id             uint      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UUID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	Results        string    `json:"results"`
	PartnerResults string    `json:"partnerResults"`
	Image          string    `json:"image"`
	PartnerImage   string    `json:"partnerImage"`
	CareOption     string    `json:"care_option"`
	UserId         uint      `json:"userId"`
	User           User      `json:"user"`
	Date           string    `json:"date"`
	Status         string    `json:"status"`
	Deleted        bool      `json:"deleted"`
}
