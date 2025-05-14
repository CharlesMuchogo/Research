package dto

import (
	"awesomeProject/models"
)

type UserDTO struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Phone          string `json:"phone"`
	Email          string `json:"email" gorm:"unique"`
	Password       string `json:"password"`
	ProfilePhoto   string `json:"profilePhoto"`
	Age            string `json:"age"`
	DeviceId       string `json:"deviceId"`
	EducationLevel string `json:"educationLevel"`
	TestedBefore   bool   `json:"testedBefore"`
	Country        string `json:"country"`
	DisplayResults bool   `json:"displayResults" gorm:"default:true"`
	Gender         string `json:"gender"`
}

func (user UserDTO) ToUser() models.User {
	return models.User{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Phone:          user.Phone,
		Email:          user.Email,
		Password:       user.Password,
		ProfilePhoto:   user.ProfilePhoto,
		Age:            user.Age,
		EducationLevel: user.EducationLevel,
		TestedBefore:   user.TestedBefore,
		Gender:         user.Gender,
		Country:        user.Country,
		Role:           "user",
	}
}
