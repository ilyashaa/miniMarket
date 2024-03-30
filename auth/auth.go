package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

func RegisterUser(c *gin.Context) {

	salt, err := generateRandomBytes(16)
	if err != nil {
		salt = []byte("test")
	}
	login := c.PostForm("login")
	password := c.PostForm("password")
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	user := User{
		Login:        login,
		HashPassword: hashedPassword,
		Salt:         salt,
	}
	Users[login] = user

	response := fmt.Sprintf("Получены данные: Login - %s , Password - %s", login, password)
	c.String(http.StatusOK, response)

}

func AuthorizationUser(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	user, ok := Users[login]

	testHashedPassword := argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32)

	if !ok || !reflect.DeepEqual(user.HashPassword, testHashedPassword) {
		c.String(http.StatusForbidden, "%s", "Неверный логин или пароль")
		return
	}

	response := fmt.Sprintf("Пароль подошёл к %s", user.Login)
	c.String(http.StatusOK, response)
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
