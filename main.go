package main

import (
	"fmt"

	"crypto/rand"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

type User struct {
	Login        string
	HashPassword []byte
	Salt         []byte
}

func main() {
	users := make(map[string]User)

	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/user", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", users)
	})

	router.POST("/user", func(c *gin.Context) {

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
		users[login] = user

		response := fmt.Sprintf("Получены данные: Login - %s , Password - %s", login, password)
		c.String(http.StatusOK, response)
	})

	router.POST("/auth", func(c *gin.Context) {

		login := c.PostForm("login")
		password := c.PostForm("password")

		user, ok := users[login]

		testHashedPassword := argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32)

		if !ok || !reflect.DeepEqual(user.HashPassword, testHashedPassword) {
			c.String(http.StatusForbidden, "%s", "Неверный логин или пароль")
			return
		}

		response := fmt.Sprintf("Пароль подошёл к %s", user.Login)
		c.String(http.StatusOK, response)
	})

	router.LoadHTMLGlob("templates/*")

	router.Run(":8080")
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
