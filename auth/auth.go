package auth

import (
	"database/sql"
	"log"
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
	Users map[string]*User
}

func NewAuth() *Auth {
	return &Auth{
		Users: make(map[string]*User),
	}
}

func (auth *Auth) Register(email string, password string, db sql.DB) string {

	validMail := isValidEmail(email)

	if !validMail {
		return "Указан некорретный email. Пройдите регистрацию заново."
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	sqlStatement := `
    INSERT INTO users (email, password)
    VALUES ($1, $2);`

	result, err := db.Exec(sqlStatement, email, passwordHash)
	if err != nil {
		log.Fatal(err)
		return "Не получилось передать данные на сервер."
	}

	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return "Вы прошли регистрацию, " + email
}

func (auth *Auth) Authorize(email string, password string, db sql.DB) string {
	query := `SELECT Email, Password FROM users WHERE email = $1`

	rows, err := db.Query(query, email)
	if err != nil {
		return "Не удалось получить данные"
	}

	defer rows.Close()
	var sqlEmail string
	var sqlPassword string

	for rows.Next() {
		err = rows.Scan(&sqlEmail, &sqlPassword)
		if err != nil {
			return "Не удалось расшифровать данные"
		}

		match, err := argon2id.ComparePasswordAndHash(password, sqlPassword)
		if err != nil {
			log.Fatal(err)
		}

		if !match {
			return "Неверный логин или пароль"
		}
	}

	return email
}

func (auth *Auth) Update(email string, nickname string, birthdayDate string) {
	user := auth.Users[email]

	if nickname != "" {
		user.Nickname = nickname
	}

	if birthdayDate != "" {
		t, err := time.Parse("2006-01-02", birthdayDate)
		if err != nil {
			panic(err)
		}
		user.BirthdayDate = t
	}
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return emailRegex.MatchString(email)
}
