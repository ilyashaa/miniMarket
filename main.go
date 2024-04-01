package main

import (
	"miniMarket/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	authService := auth.NewAuth()

	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/user", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s")
	})

	router.POST("/user", func(c *gin.Context) {
		login := c.PostForm("login")
		password := c.PostForm("password")
		response := authService.Register(login, password)
		c.String(http.StatusOK, response)
	})

	router.POST("/auth", func(c *gin.Context) {
		login := c.PostForm("login")
		password := c.PostForm("password")
		response := authService.Authorize(login, password)
		if response.RequestError {
			c.String(http.StatusOK, response.Text)
		}
		c.String(http.StatusForbidden, response.Text)
	})

	router.LoadHTMLGlob("templates/*")

	router.Run(":8080")
}
