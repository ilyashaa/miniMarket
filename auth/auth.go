package auth

import (
	"crypto/rand"

	"reflect"

	"golang.org/x/crypto/argon2"
)

type User struct {
	Login        string
	HashPassword []byte
	Salt         []byte
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

func (auth *Auth) Register(login string, password string) string {

	salt, err := generateRandomBytes(16)
	if err != nil {
		salt = []byte("test")
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	user := User{
		Login:        login,
		HashPassword: hashedPassword,
		Salt:         salt,
	}
	auth.Users[login] = user

	return login
}

func (auth *Auth) Authorize(login string, password string) string {
	user, ok := auth.Users[login]
	if !ok {
		return "*неверный логин или пароль*"
	}

	testHashedPassword := argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32)

	if !reflect.DeepEqual(user.HashPassword, testHashedPassword) {
		return "*неверный логин или пароль*"
	}

	return login
}

// func (auth *Auth) CheckList(login string, password string) string {
// 	if _, ok := auth.Users[login]; ok {
// 		return auth.Authorize(login, password)
// 	} else {
// 		// Вместо регистрации нового пользователя возвращаем сообщение об ошибке
// 		return "*неверный логин или пароль*"
// 	}
// }

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
