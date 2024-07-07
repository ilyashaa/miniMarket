package userDB

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     string
	Nickname     string    `gorm:"type:varchar(15); default:Nickname"`
	BirthdayDate time.Time `gorm:"type:date"`
}

func RegisterSQL(email string, passwordHash string, db *gorm.DB) (string, error) {

	user := User{Email: email, Password: passwordHash}
	result := db.Create(&user)
	if result.Error != nil {
		return "Не получилось передать данные на сервер.", nil
	}

	return "Вы прошли регистрацию, " + email, nil
}

func AuthorizeSQL(email string, password string, db *gorm.DB) (User, error) {

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil

}

func SelectInfoSQL(IdUser string, db *gorm.DB) (User, error) {

	var user User
	result := db.Where("id = ?", IdUser).First(&user)
	if result.Error != nil {
		return User{}, nil
	}

	return user, nil
}

func UpdateInfoSQL(email, nickname, birthdayDate string, db *gorm.DB) {

	birthDate, err := time.Parse(time.DateOnly, birthdayDate)
	if err != nil {
		log.Fatal(err)
	}

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Fatal(err)
	}

	updates := User{
		Nickname:     nickname,
		BirthdayDate: birthDate,
	}
	db.Model(&user).Updates(updates)
}
