package productDB

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id    int
	Name  string
	Price float32
	Image string
}

func LocalProduct(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price, image FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Image); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}


func GetProducts(db *sql.DB, ids []int) {
	// получения списка выбранных товаров, возвращаю id товара и его цену. Работает с массивом id
}

func CreateOrder(db *sql.DB, c *gin.Context) {
	var selectedItems map[int]int

	if err := c.BindJSON(&selectedItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unmarshal JSON"})
		return
	}

}
