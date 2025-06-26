package database

import "awesomeProject/models"

func GetUserById(id uint) (models.User, error) {
	var user models.User

	if err := DbInstance.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
