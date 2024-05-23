package saledb

import (
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Sale struct {
	Id       string          `gorm:"primaryKey"`
	AllCost  decimal.Decimal `gorm:"type:decimal"`
	SaleTime time.Time       `gorm:"type:date"`
}

type ProductsInSale struct {
	IdSale       string // ЗАМЕНИТЬ ВЕЗДЕ UINT НА INT
	IdProduct    string // gorm записи
	CostProduct  decimal.Decimal
	CountProduct int
}

func CreateSale(cost decimal.Decimal, selectedProducts map[int]int, db *gorm.DB) (string, error) { // + возвращать id, error
	saleTime := time.Now().UTC()
	saleId, err := gonanoid.New()
	if err != nil {
		return "0", err // доработать
	}
	sale := Sale{
		Id:       saleId,
		AllCost:  cost,
		SaleTime: saleTime,
	}
	result := db.Create(&sale)
	if result.Error != nil {
		fmt.Printf("error: %v\n", result.Error)
	}
	return sale.Id, nil
	// var saleID Sale
	// resultSaleID := db.Where("SaleTime = ?", sale.SaleTime).First(&saleID)
	// if resultSaleID.Error != nil {
	// 	fmt.Printf("error: %v\n", resultSaleID.Error)
	// }
}

// func AddProductsToSale(saleID string, selectedProducts map[int]int, db *gorm.DB) {
// 	priceProduct := productDB.GetPriceProduct(db)
// }
