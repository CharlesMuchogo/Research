package models

type Results struct {
	Id             uint64 `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Results        string `json:"results"`
	PartnerResults string `json:"partnerResults"`
	Image          string `json:"image"`
	PartnerImage   string `json:"partnerImage"`
	CareOption     string `json:"care_option"`
	UserId         uint   `json:"userId"`
	User           User   `json:"user"`
	Date           string `json:"date"`
}
