package main

import (
	"fmt"
	"log"

	"miniMarket/auth"
	"miniMarket/db/productDB"
	product "miniMarket/db/productDB"

	saleDB "miniMarket/db/saleDB"
	userDB "miniMarket/db/userDB"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	db := StartDB()

	defer CloseDB(db)

	db.AutoMigrate(&userDB.User{}, &productDB.Product{}, &saleDB.Sale{}, &saleDB.ProductsInSale{})

	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/products", func(c *gin.Context) {
		products, err := product.LocalProduct(db)
		if err != nil {
			return
		}
		formattedProducts := make([]gin.H, len(products))
		for i, product := range products {
			formattedProducts[i] = gin.H{
				"Id":    product.Id,
				"Name":  product.Name,
				"Price": product.Price.StringFixed(2),
				"Image": product.Image,
			}
		}
		c.HTML(http.StatusOK, "products.html", gin.H{
			"products": formattedProducts,
		})
	})

	router.POST("/products", func(c *gin.Context) {
		// GetProduct(productDB), CostProduct(тут же), CreateSale(saleDB), AddProductsToSale(saleDB)
		// product.CreateOrder(db, c) // пенести в orderDB
		// достать id продуктов из запроса http
		selectedProduct, idProducts := productDB.GetProduct(db, c)
		costProducts, err := productDB.GetPriceProduct(db, idProducts)
		if err != nil {
			log.Fatal(err)
		}
		var allCostProducts float64
		for id, count := range selectedProduct {
			allCostProducts = allCostProducts + (float64(count) * costProducts[id-1])
		}
		saleID, err := saleDB.CreateSale(allCostProducts, selectedProduct, db)
		if err != nil {
			log.Fatal(err)
		}
		saleDB.AddProductsToSale(saleID, selectedProduct, costProducts, db)

		c.Redirect(http.StatusSeeOther, "http://localhost:8080/product")
	})

	router.GET("/home", func(c *gin.Context) {

		c.HTML(http.StatusOK, "home.html", gin.H{
			"Email": "пользователь",
		})
	})

	router.POST("/hello", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		result := auth.Register(email, password, db)
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
		result := auth.Authorize(email, password, db)
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

		sqlEmail, sqlNickname, sqlBirthdayDate, err := userDB.SelectInfoSQL(emailKey, db)
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "user.html", gin.H{
			"Email":        sqlEmail,
			"Nickname":     sqlNickname,
			"BirthdayDate": sqlBirthdayDate,
		})
	})

	router.POST("/user", func(c *gin.Context) {
		email := c.PostForm("email")
		nickname := c.PostForm("nickname")
		birthDate := c.PostForm("birthdayDate")
		userDB.UpdateInfoSQL(email, nickname, birthDate, db)
		c.Redirect(http.StatusSeeOther, "http://localhost:8080/user")
	})

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/home")
	})

	router.Run(":8080")
}

func StartDB() *gorm.DB {
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

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}
