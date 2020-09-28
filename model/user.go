package model

import "gorm.io/gorm"

// Model User
type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func CreateUser(data User) error {
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
