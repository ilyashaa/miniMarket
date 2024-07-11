package auth

import (
	"log"
	store "miniMarket/db/userDB"
	"regexp"
	"time"

	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

type User struct {
	Email        string
	HashPassword string
	Nickname     string
	BirthdayDate time.Time
}

type Response struct {
	RequestError bool
	Text         string
}

func Register(email string, password string, db *gorm.DB) string {

	validMail := isValidEmail(email)

	if !validMail {
		return "Указан некорретный email. Пройдите регистрацию заново."
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	result, err := store.RegisterSQL(email, passwordHash, db)
	if err != nil {
		log.Fatal(err)
		return result
	}

	return result
}

func Authorize(email, password string, db *gorm.DB) store.User {

	user, err := store.AuthorizeSQL(email, password, db)
	if err != nil {
		log.Panicln(err)
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.Password)
	if err != nil {
		log.Panicln(err)
	}

	if !match {
		return store.User{}
	}

	return user
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return emailRegex.MatchString(email)
}
