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

// Model Login
type UserLogin struct {
	Email    string
	Password string
}

func CreateUser(data User) error {
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func FindByEmail(email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &user, nil
}

func FindById(id string) (*User, error) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &user, nil
}

func DeleteById(id string) error {
	err := db.Delete(&User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser() (*[]User, error) {
	var user []User
	err := db.Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateByUser(user *User) error {
	err := db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}
