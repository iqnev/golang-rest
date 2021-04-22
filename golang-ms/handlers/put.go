package handlers

import (
	"net/http"

	"github.com/iqnev/golang-rest/ms/data"
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
	rw.Header().Add("Content-Type", "application/json")
	product := req.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("Updating record id", product.ID)

	dataErr := p.productDB.UpdateProduct(&product)

	if dataErr == data.ErrProductNotFound {
		p.l.Error("Product not found", dataErr)
		rw.WriteHeader(http.StatusNotFound)

		data.ToJson(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}
