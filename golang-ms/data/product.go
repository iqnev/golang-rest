package data

import (
	"fmt"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"  validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price"  validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeleteOn    string  `json:"-"`
}

type Products []*Product

var ErrProductNotFound = fmt.Errorf("Product not found")

func GetProducts() Products {
	return productList
}

func UpdateProduct(id int, pr *Product) error {
	_, pos, err := findProduct(id)

	if err != nil {
		return err
	}

	pr.ID = id
	productList[pos] = pr

	return nil
}

func DeleteProduct(id int) error {
	i := findIndexByProductID(id)

	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

func findIndexByProductID(id int) int {

	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound

}

func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

func AddProduct(pr *Product) {
	pr.ID = getNextID()

	productList = append(productList, pr)
}

func getNextID() int {
	lp := productList[len(productList)-1]

	return lp.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Strong milk coffee",
		Price:       1.99,
		SKU:         "def123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
