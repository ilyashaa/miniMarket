package product

import "miniMarket/db/store"

type Product struct {
	Name  string
	Price float32
	Image string
}

func LocalProduct() ([]Product, error) {
	result, err := store.GetProduct()
	return result, err
}
