package auth

import (
	"my-project/configs/database"
)

func CreateUser(req *RegisterRequest) error {
	user := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	return database.DB.Create(&user).Error
}

func GetUsers() ([]User, error) {
	var users []User
	err := database.DB.Find(&users).Error
	return users, err
}

func GetUserByID(id uint) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func DeleteUser(id uint) error {
	return database.DB.Delete(&User{}, id).Error
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err

}
