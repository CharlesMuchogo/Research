package models

type Message struct {
	Id        uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp" gorm:"not null;index:idx_user_timestamp,unique"`
	UserId    uint   `json:"-" gorm:"not null;index:idx_user_timestamp,unique"`
}
