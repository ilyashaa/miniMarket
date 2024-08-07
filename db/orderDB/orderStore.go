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
	IdUser    uint      `gorm:"type:int"`
	TimeOrder time.Time `gorm:"type:date"`
	Status    string    `gorm:"type:text"`
}

type OrderSaleInfoSQL struct {
	IdOrder      string
	IdSale       string
	IdUser       uint
	TimeOrder    time.Time
	Status       string
	AllCost      decimal.Decimal
	NameProduct  string
	CostProduct  decimal.Decimal
	CountProduct int
}

type OrderSaleInfo struct {
	IdOrder     string
	IdSale      string
	IdUser      uint
	TimeOrder   time.Time
	Status      string
	AllCost     decimal.Decimal
	ProductList []ProductList
}

type ProductList struct {
	NameProduct  string
	CostProduct  decimal.Decimal
	CountProduct int
}

func CreateOrder(saleID string, IdUser string, db *gorm.DB) Order {
	orderTime := time.Now().UTC()
	orderID, err := gonanoid.New()
	if err != nil {
		log.Fatal(err)

	}
	var user userDB.User
	db.Where("id = ?", IdUser).First(&user)
	order := Order{
		IdOrder:   orderID,
		IdSale:    saleID,
		IdUser:    user.ID,
		TimeOrder: orderTime,
		Status:    "Ожидает оплаты",
	}

	db.Create(&order)
	return order
}

func GetOrders(db *gorm.DB, IdUser string) []Order {

	var orders []Order
	db.Where("id_user = ?", IdUser).Find(&orders)
	return orders
}

func GetInfoOrder(idOrder string, db *gorm.DB) OrderSaleInfo {
	var orders []OrderSaleInfoSQL

	db.Table("orders").Select("orders.id_order", "orders.id_sale", "orders.id_user", "orders.time_order", "orders.status",
		"sales.all_cost",
		"products_in_sales.cost_product", "products_in_sales.count_product",
		"products.name_product").
		Joins("JOIN sales ON sales.id = orders.id_sale").
		Joins("JOIN products_in_sales ON products_in_sales.id_sale = sales.id").
		Joins("JOIN products ON products.id = products_in_sales.id_product").
		Where("orders.id_order = ?", idOrder).
		Find(&orders)

	order := OrderSaleInfo{
		IdOrder:   orders[0].IdOrder,
		IdSale:    orders[0].IdSale,
		IdUser:    orders[0].IdUser,
		TimeOrder: orders[0].TimeOrder,
		Status:    orders[0].Status,
		AllCost:   orders[0].AllCost,
	}

	for i := range orders {
		order.ProductList = append(order.ProductList, ProductList{
			NameProduct:  orders[i].NameProduct,
			CostProduct:  orders[i].CostProduct,
			CountProduct: orders[i].CountProduct,
		})
	}
	return order
}

func UpdateInfoOrder(idOrder string, db *gorm.DB) OrderSaleInfo {
	order := GetInfoOrder(idOrder, db)
	order.Status = "Оплачен"
	var orderSQL Order
	db.Where("id_order = ?", idOrder).First(&orderSQL)
	orderSQL.Status = "Оплачен"
	db.Save(&orderSQL)
	return order
}
