package repository

import "gorm.io/gorm"

func FindUserByName(username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(username string, password string) (*User, error) {
	user := User{
		Username: username,
		Password: password,
	}
	err := db.Create(&user).Error

	return &user, err

}
