package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
Name string
Email string
Age int
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
		ctx.String(http.StatusOK, "%s",users)
	})

	router.POST("/user", func(c *gin.Context) {
		
		name := c.PostForm("name")
		email := c.PostForm("email")
		age, _ := strconv.Atoi(c.PostForm("age"))

		user := User{
			Name: name,
			Email: email,
			Age: age,
		}
        users[email] = user

		response := fmt.Sprintf("Получены данные: Имя - %s , Email - %s, Age - %d ", name, email, age )
		c.String(http.StatusOK, response)
	})

	router.POST("/auth", func(c *gin.Context) {
		
		name := c.PostForm("name")
		email := c.PostForm("email")

        user, ok := users[email]

		if !ok || strings.Compare(user.Name, name) != 0 {
			c.String(http.StatusForbidden, "%s","Неверный логин или пароль")
			return
		}

		response := fmt.Sprintf("Получены данные: Имя - %s , Email - %s, Age - %d ", user.Name, user.Email, user.Age)
		c.String(http.StatusOK, response)
	})
	

	router.LoadHTMLGlob("templates/*")

	router.Run(":8080")
}
