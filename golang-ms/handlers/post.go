package handlers

import (
	"net/http"

	"github.com/iqnev/golang-rest/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *Products) Create(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST Product")

	product := req.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
}
