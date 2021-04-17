package handlers

import (
	"net/http"

	"github.com/iqnev/golang-rest/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, req *http.Request) {
	product := req.Context().Value(KeyProduct{}).(data.Product)
	p.l.Println("[DEBUG] updating record id", product.ID)

	dataErr := data.UpdateProduct(product.ID, &product)

	if dataErr == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", dataErr)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}
