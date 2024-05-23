package productDB

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	Id    string          `gorm:"primaryKey"`
	Name  string          `gorm:"type:varchar(10)"`
	Price decimal.Decimal `gorm:"type:decimal"`
	Image string          `gorm:"type:text"`
}

func LocalProduct(db *gorm.DB) ([]Product, error) {
	var products []Product
	db.Find(&products)

	return products, nil
}

func GetPriceProduct(db *gorm.DB, idsProducts []int) map[int]float64 { // + возвращать  error
	idAndPriceProduct := make(map[int]float64)
	for _, id := range idsProducts {
		var product Product
		// SQL: "Select * FROM products WHERE id in(id1,id2...)" - перевести в Go
		result := db.Where("id in ?", idsProducts).First(&product)
		if result.Error != nil {
			return nil
		}
		idAndPriceProduct[id] = product.Price
	}

	return idAndPriceProduct
}

// func GetProduct(db *gorm.DB, c *gin.Context) (map[int]int, []int) {
// 	var selectedItems map[int]int

// 	if err := c.BindJSON(&selectedItems); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unmarshal JSON"})
// 		log.Fatal(err)
// 	}
// 	var idsProducts []int
// 	for idProducts := range selectedItems {
// 		idsProducts = append(idsProducts, idProducts)
// 	}
// 	return selectedItems, idsProducts
// }
