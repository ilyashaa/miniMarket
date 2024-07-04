package saledb

import (
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Sale struct {
	gorm.Model
	Id       string          `gorm:"primaryKey"`
	AllCost  decimal.Decimal `gorm:"type:decimal"`
	SaleTime time.Time       `gorm:"type:date"`
}

type ProductsInSale struct {
	gorm.Model
	IdSale       string          `gorm:"type:text"`
	IdProduct    int             `gorm:"type:integer"`
	CostProduct  decimal.Decimal `gorm:"type:decimal"`
	CountProduct int             `gorm:"type:integer"`
}

func CreateSale(cost float64, selectedProducts map[int]int, db *gorm.DB) (string, error) { // + возвращать id, error
	saleTime := time.Now().UTC()
	saleId, err := gonanoid.New()
	if err != nil {
		return "0", err // доработать
	}
	allCost := decimal.NewFromFloat(cost)
	sale := Sale{
		Id:       saleId,
		AllCost:  allCost,
		SaleTime: saleTime,
	}
	result := db.Create(&sale)
	if result.Error != nil {
		fmt.Printf("error: %v\n", result.Error)
	}
	return sale.Id, nil
}

func AddProductsToSale(saleID string, selectedProducts map[int]int, idAndPrice map[int]float64, db *gorm.DB) {
	for idProduct, countProduct := range selectedProducts {
		productInSale := ProductsInSale{
			IdSale:       saleID,
			IdProduct:    idProduct,
			CostProduct:  decimal.NewFromFloat(idAndPrice[idProduct]),
			CountProduct: countProduct,
		}
		result := db.Create(&productInSale)
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
		}
	}
}
