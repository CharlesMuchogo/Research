package models

type Clinic struct {
	ID       uint   `gorm:"primaryKey;autoIncrement:true;unique" json:"id"`
	Name     string `gorm:"unique" json:"name"`
	Address  string `json:"address"`
	Contacts string `json:"contacts"`
	Active   bool   `json:"active"`
}
