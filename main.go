package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.POST("/user", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		age := c.PostForm("age")
		response := "Получены данные: Имя - " + name + ", Email - " + email + ", Age - " + age
		c.String(http.StatusOK, response)
	})

	router.LoadHTMLFiles("templates/author.html", "templates/register.html")

	router.Run(":8080")
}
