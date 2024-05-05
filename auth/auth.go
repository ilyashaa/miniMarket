package auth

import (
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

type Auth struct {
	UserStore *store.UserStore
}

func NewAuth() *Auth {
	return &Auth{
		UserStore: store.NewUserStore(),
	}
}

func (auth *Auth) Register(email string, password string) string {

	validMail := isValidEmail(email)

	if !validMail {
		return "Указан некорретный email. Пройдите регистрацию заново."
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	result, err := auth.UserStore.RegisterSQL(email, passwordHash)
	if err != nil {
		log.Fatal(err)
		return result
	}

	return result
}

func (auth *Auth) Authorize(email, password string) string {

	sqlPassword, err := auth.UserStore.AuthorizeSQL(email, password)
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

func (auth *Auth) Update(email, nickname, birthdayDate string) {
	// user := auth.Users[email]

	// if nickname != "" {
	// 	user.Nickname = nickname
	// }

	// if birthdayDate != "" {
	// 	t, err := time.Parse("2006-01-02", birthdayDate)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	user.BirthdayDate = t
	// }
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return emailRegex.MatchString(email)
}
