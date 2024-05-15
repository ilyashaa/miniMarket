package auth

import (
	"database/sql"
	"log"
	"miniMarket/db/store"
	"regexp"
	"time"

	"github.com/alexedwards/argon2id"
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

func Register(email string, password string, db *sql.DB) string {

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

func Authorize(email, password string, db *sql.DB) string {

	sqlPassword, err := store.AuthorizeSQL(email, password, db)
	if err != nil {
		log.Fatal(err)
	}

	match, err := argon2id.ComparePasswordAndHash(password, sqlPassword)
	if err != nil {
		log.Fatal(err)
	}

	if !match {
		return "Неверный логин или пароль"
	}

	return email
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return emailRegex.MatchString(email)
}
