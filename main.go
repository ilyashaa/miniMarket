package main

import (
	"miniMarket/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/user", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", auth.Users)
	})

	router.POST("/user", func(c *gin.Context) {
		auth.RegisterUser(c)
	})

	router.POST("/auth", func(c *gin.Context) {
		auth.AuthorizationUser(c)
	})

	router.LoadHTMLGlob("templates/*")

	router.Run(":8080")
}
