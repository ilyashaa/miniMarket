package product

type Product struct {
	Name  string
	Price float32
	Image string
}

type ListProduct struct {
	Products []Product
}

func NewProduct() *ListProduct {
	products := &ListProduct{
		Products: []Product{},
	}
	products.initProducts()
	return products
}

func (lp *ListProduct) initProducts() {
	lp.AddProduct(Product{Name: "Яблоко", Price: 27.8, Image: "static/images/apple.svg"})
	lp.AddProduct(Product{Name: "Банан", Price: 23.5, Image: "static/images/banana.svg"})
	lp.AddProduct(Product{Name: "Виноград", Price: 23.5, Image: "static/images/grapes.svg"})
	lp.AddProduct(Product{Name: "Апельсин", Price: 23.5, Image: "static/images/orange.svg"})
	lp.AddProduct(Product{Name: "Груша", Price: 23.5, Image: "static/images/peach.svg"})
	lp.AddProduct(Product{Name: "Арбуз", Price: 23.5, Image: "static/images/watermelon.svg"})
}

func (lp *ListProduct) AddProduct(p Product) {
	lp.Products = append(lp.Products, p)
}

func (lp *ListProduct) LocalList() []Product {
	return lp.Products
}
