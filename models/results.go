package models

import (
	"gorm.io/gorm"
)

type Results struct {
	gorm.Model
	Results        string `json:"results"`
	PartnerResults string `json:"partnerResults"`
	Image          string `json:"image"`
	PartnerImage   string `json:"partnerImage"`
	CareOption     string `json:"care_option"`
}
