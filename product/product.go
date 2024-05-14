package product

import (
	"database/sql"
)

type Product struct {
	Name  string
	Price float32
	Image string
}

func LocalProduct(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT name, price, image FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Name, &p.Price, &p.Image); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
