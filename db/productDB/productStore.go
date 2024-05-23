package productDB

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Product struct {
	Id    int     `gorm:"primaryKey"`
	Name  string  `gorm:"type:varchar(10)"`
	Price float64 `gorm:"type:decimal"`
	Image string  `gorm:"type:text"`
}

func LocalProduct(db *gorm.DB) ([]Product, error) {
	var products []Product
	db.Find(&products)

	return products, nil
}

func GetProducts(db *gorm.DB, ids []int) {
	// получения списка выбранных товаров, возвращаю id товара и его цену. Работает с массивом id
}

func CreateOrder(db *gorm.DB, c *gin.Context) {
	var selectedItems map[int]int

	if err := c.BindJSON(&selectedItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unmarshal JSON"})
		return
	}

}
