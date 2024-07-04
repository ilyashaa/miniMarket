package main

import (
	"fmt"
	"log"
	"strconv"

	"miniMarket/auth"
	orderDB "miniMarket/db/orderDB"
	productDB "miniMarket/db/productDB"
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

	db.AutoMigrate(&userDB.User{}, &productDB.Product{}, &saleDB.Sale{}, &saleDB.ProductsInSale{}, &orderDB.Order{})

	router := gin.Default()

	router.GET("/author", func(c *gin.Context) {
		c.HTML(http.StatusOK, "author.html", gin.H{})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/products", func(c *gin.Context) {
		IdUser, err := c.Cookie("IdUser")
		if err != nil {
			IdUser = ""
			return
		}
		orders := orderDB.GetOrders(db, IdUser)
		products, err := productDB.LocalProduct(db)
		if err != nil {
			return
		}
		formattedProducts := make([]gin.H, len(products))
		for i, product := range products {
			formattedProducts[i] = gin.H{
				"Id":    product.ID,
				"Name":  product.NameProduct,
				"Price": product.Price.StringFixed(2),
				"Image": product.Image,
			}
		}
		c.HTML(http.StatusOK, "products.html", gin.H{
			"products": formattedProducts,
			"orders":   orders,
		})
	})

	router.POST("/products", func(c *gin.Context) {
		var selectedItems map[int]int

		if err := c.BindJSON(&selectedItems); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unmarshal JSON"})
			log.Fatal(err)
		}

		var idProducts []int
		for id := range selectedItems {
			idProducts = append(idProducts, id)
		}

		costProducts, err := productDB.GetPriceProduct(db, idProducts)
		if err != nil {
			log.Fatal(err)
		}
		idAndPrice := make(map[int]float64)
		var allCostProducts float64
		for idCost, price := range costProducts {
			for idProduct, count := range selectedItems {
				if idCost == idProduct {
					allCostProducts += (float64(count) * price)
					idAndPrice[idProduct] = price
				}
			}
		}
		saleID, err := saleDB.CreateSale(allCostProducts, selectedItems, db)
		if err != nil {
			log.Fatal(err)
		}
		saleDB.AddProductsToSale(saleID, selectedItems, idAndPrice, db)
		IdUser, err := c.Cookie("IdUser")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		orderDB.CreateOrder(saleID, IdUser, db)
		c.Redirect(http.StatusSeeOther, "http://localhost:8080/product")
	})

	router.GET("/order/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		order := orderDB.GetInfoOrder(orderID, db)
		c.HTML(http.StatusOK, "order.html", gin.H{
			"order": order,
		})
	})

	router.GET("/home", func(c *gin.Context) {
		email, err := c.Cookie("email")
		if err != nil {
			email = "пользователь"
		}
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Email": email,
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
		result := auth.Authorize(email, password, db)
		cookieEmail := http.Cookie{
			Name:    "IdUser",
			Value:   strconv.FormatUint(uint64(result.ID), 10),
			Expires: expiration,
		}
		http.SetCookie(c.Writer, &cookieEmail)
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Email": result.Email,
		})

	})

	router.GET("/user", func(c *gin.Context) {
		IdUser, err := c.Cookie("IdUser")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user, err := userDB.SelectInfoSQL(IdUser, db)
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "user.html", gin.H{
			"user": user,
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
