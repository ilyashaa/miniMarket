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

	router.GET("/product", func(c *gin.Context) {
		c.HTML(http.StatusOK, "product.html", gin.H{})
	})

	router.GET("/user", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", gin.H{
			"Login": "пользователь",
		})
	})

	router.POST("/hello", func(c *gin.Context) {
		login := c.PostForm("login")
		password := c.PostForm("password")
		result := authService.Register(login, password)
		c.HTML(http.StatusOK, "hello.html", gin.H{
			"Login": result,
		})
	})

	router.POST("/user", func(c *gin.Context) {
		login := c.PostForm("login")
		password := c.PostForm("password")
		result := authService.Authorize(login, password)
		c.HTML(http.StatusOK, "user.html", gin.H{
			"Login": result,
		})
	})

	// router.POST("/auth", func(c *gin.Context) {
	// 	login := c.PostForm("login")
	// 	password := c.PostForm("password")
	// 	response := authService.Authorize(login, password)

	// 	if response.RequestError {
	// 		c.String(http.StatusOK, response.Text)
	// 	}
	// 	c.String(http.StatusForbidden, response.Text)
	// })

	router.LoadHTMLGlob("templates/*")

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/user")
	})

	router.Run(":8080")
}
