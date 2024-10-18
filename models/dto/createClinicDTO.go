package dto

import "awesomeProject/models"

type CreateClinicDTO struct {
	Name        string `gorm:"unique" json:"name"`
	Address     string `json:"address"`
	Contacts    string `json:"contacts"`
	Active      bool   `json:"active"`
	Description string `json:"description"`
}

func (c CreateClinicDTO) ToClinic() models.Clinic {
	return models.Clinic{
		Name:        c.Name,
		Address:     c.Address,
		Contacts:    c.Contacts,
		Active:      true,
		Description: c.Description,
	}
}
