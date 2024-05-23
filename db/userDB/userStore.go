package userDB

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint   `gorm:"type:uuid; primaryKey"`
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

func AuthorizeSQL(email string, password string, db *gorm.DB) (string, error) {

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "Не удалось получить данные", nil
	}

	return user.Password, nil

}

func SelectInfoSQL(emailKey string, db *gorm.DB) (string, string, time.Time, error) {

	var user User
	result := db.Where("email = ?", emailKey).First(&user)
	if result.Error != nil {
		defaultTime := time.Now()
		return "Не удалось получить данные", "Не удалось получить данные", defaultTime, nil
	}

	return user.Email, user.Nickname, user.BirthdayDate, nil
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
