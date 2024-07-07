package productDB

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	NameProduct string          `gorm:"type:varchar(10)"`
	Price       decimal.Decimal `gorm:"type:decimal"`
	Image       string          `gorm:"type:text"`
}

func LocalProduct(db *gorm.DB) ([]Product, error) {
	var products []Product
	db.Find(&products)

	return products, nil
}

func GetPriceProduct(db *gorm.DB, idsProducts []int) (map[int]float64, error) {
	priceProduct := make(map[int]float64)
	var products []Product
	result := db.Select("id", "price").Where("id IN ?", idsProducts).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, product := range products {
		price, _ := product.Price.Float64()
		priceProduct[int(product.ID)] = price
	}

	return priceProduct, nil
}
