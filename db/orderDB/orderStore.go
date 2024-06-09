package orderdb

import (
	"log"
	"miniMarket/db/userDB"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/shopspring/decimal"
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

type OrderSaleInfoSQL struct {
	IdOrder      string
	IdSale       string
	EmailUser    string
	TimeOrder    time.Time
	Status       string
	AllCost      decimal.Decimal
	NameProduct  string
	CostProduct  decimal.Decimal
	CountProduct int
}

type OrderSaleInfo struct {
	IdOrder      string
	IdSale       string
	EmailUser    string
	TimeOrder    time.Time
	Status       string
	AllCost      decimal.Decimal
	NameProduct  []string
	CostProduct  []decimal.Decimal
	CountProduct []int
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

func GetInfoOrder(idOrder string, db *gorm.DB) OrderSaleInfo {
	var orders []OrderSaleInfoSQL

	db.Table("orders").Select("orders.id_order", "orders.id_sale", "orders.email_user", "orders.time_order", "orders.status",
		"sales.all_cost",
		"products_in_sales.cost_product", "products_in_sales.count_product",
		"products.name_product").
		Joins("JOIN sales ON sales.id = orders.id_sale").
		Joins("JOIN products_in_sales ON products_in_sales.id_sale = sales.id").
		Joins("JOIN products ON products.id = products_in_sales.id_product").
		Where("orders.id_order = ?", idOrder).
		Find(&orders)
		// db.Where("id_order = ?", idOrder).First(&order)

	order := OrderSaleInfo{
		IdOrder:   orders[0].IdOrder,
		IdSale:    orders[0].IdSale,
		EmailUser: orders[0].EmailUser,
		TimeOrder: orders[0].TimeOrder,
		Status:    orders[0].Status,
		AllCost:   orders[0].AllCost,
	}

	for i := range orders {
		order.NameProduct = append(order.NameProduct, orders[i].NameProduct)
		order.CostProduct = append(order.CostProduct, orders[i].CostProduct)
		order.CountProduct = append(order.CountProduct, orders[i].CountProduct)
	}
	return order
}
