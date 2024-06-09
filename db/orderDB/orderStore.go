package orderdb

import (
	"log"
	"miniMarket/db/userDB"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	IdOrder   string    `gorm:"primaryKey"`
	IdSale    string    `gorm:"type:text"`
	EmailUser string    `gorm:"type:text"`
	TimeOrder time.Time `gorm:"type:date"`
	Status    string    `gorm:"type:text"`
}

func CreateOrder(saleID string, emailKey string, db *gorm.DB) {
	orderTime := time.Now().UTC()
	orderID, err := gonanoid.New()
	if err != nil {
		log.Fatal(err)
	}
	var user userDB.User
	result := db.Where("email = ?", emailKey).First(&user)
	if result.Error != nil {
		log.Fatal(err)
	}
	order := Order{
		IdOrder:   orderID,
		IdSale:    saleID,
		EmailUser: user.Email,
		TimeOrder: orderTime,
		Status:    "Pending",
	}
	db.Create(&order)
}

func GetOrders(db *gorm.DB) []Order {
	var orders []Order
	db.Find(&orders)
	return orders
}

func GetInfoOrder(idOrder string, db *gorm.DB) Order {
	var order Order
	db.Where("id_order = ?", idOrder).First(&order)
	return order
}
