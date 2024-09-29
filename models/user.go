package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint   `gorm:"primaryKey;autoIncrement:true;unique" json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Phone          string `json:"phone"`
	Email          string `json:"email" gorm:"unique"`
	Password       string `json:"-"`
	DeviceId       string `json:"-"`
	ProfilePhoto   string `json:"profilePhoto"`
	Age            string `json:"age"`
	EducationLevel string `json:"educationLevel"`
	TestedBefore   bool   `json:"testedBefore"`
	SaveResults    bool   `json:"SaveResults"`
	Gender         string `json:"gender"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
