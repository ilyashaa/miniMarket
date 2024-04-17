package auth

import (
	"crypto/rand"
	"time"

	"reflect"
	"regexp"

	"golang.org/x/crypto/argon2"
)

type User struct {
	Email        string
	HashPassword []byte
	Salt         []byte
	Nickname     *string
	BirthdayDate *time.Time
}

type Response struct {
	RequestError bool
	Text         string
}

type Auth struct {
	Users map[string]User
}

func NewAuth() *Auth {
	return &Auth{
		Users: make(map[string]User),
	}
}

func (auth *Auth) Register(email string, password string) string {

	validMail := isValidEmail(email)

	if !validMail {
		return "Указан некорретный email. Пройдите регистрацию заново."
	}

	salt, err := generateRandomBytes(16)
	if err != nil {
		salt = []byte("test")
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	user := User{
		Email:        email,
		HashPassword: hashedPassword,
		Salt:         salt,
	}
	auth.Users[email] = user

	return "Вы прошли регистрацию, " + email
}

func (auth *Auth) Authorize(email string, password string) string {
	user, ok := auth.Users[email]
	if !ok {
		return "*неверный логин или пароль*"
	}

	testHashedPassword := argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32)

	if !reflect.DeepEqual(user.HashPassword, testHashedPassword) {
		return "*неверный логин или пароль*"
	}

	return email
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	return emailRegex.MatchString(email)
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
