package main

import (
	"database/sql"
	"fmt"
	"log"
	"miniMarket/auth"
	"miniMarket/product"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "admin"
		password = "12345"
		dbname   = "postgres"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewAuth()

	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/products", func(c *gin.Context) {
		products, err := product.LocalProduct(*db)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "products.html", gin.H{
			"products": products,
		})
	})

	router.GET("/home", func(c *gin.Context) {

		c.HTML(http.StatusOK, "home.html", gin.H{
			"Email": "пользователь",
		})
	})

	router.POST("/hello", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		result := authService.Register(email, password, *db)
		c.HTML(http.StatusOK, "hello.html", gin.H{
			"Email": result,
		})
	})

	router.POST("/home", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:    "email",
			Value:   email,
			Expires: expiration,
		}
		http.SetCookie(c.Writer, &cookie)
		result := authService.Authorize(email, password, *db)
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Email": result,
		})

	})

	router.GET("/user", func(c *gin.Context) {
		emailKey, err := c.Cookie("email")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		query := `SELECT email, nickname, birthdaydate FROM users WHERE email = $1`

		rows, err := db.Query(query, emailKey)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		defer rows.Close()
		var sqlEmail, sqlNickname, sqlBirthdayDate *string

		for rows.Next() {
			err = rows.Scan(&sqlEmail, &sqlNickname, &sqlBirthdayDate)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

		}
		if sqlBirthdayDate == nil {
			c.HTML(http.StatusOK, "user.html", gin.H{
				"Email":    sqlEmail,
				"Nickname": sqlNickname,
			})
			return
		}
		birthdayDate, err := time.Parse("2006-01-02", *sqlBirthdayDate)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "user.html", gin.H{
			"Email":        sqlEmail,
			"Nickname":     sqlNickname,
			"BirthdayDate": birthdayDate,
		})
	})

	router.POST("/user", func(c *gin.Context) {
		email := c.PostForm("email")
		nickname := c.PostForm("nickname")
		birthDate := c.PostForm("birthDate")
		authService.Update(email, nickname, birthDate)
		c.Redirect(http.StatusSeeOther, "http://localhost:8080/user")
	})

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/home")
	})

	router.Run(":8080")
}
