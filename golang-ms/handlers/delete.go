package handlers

import (
	"net/http"

	"github.com/iqnev/golang-rest/ms/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *Products) Delete(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(req)

	p.l.Debug("Deleting record", "id", id)

	err := p.productDB.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		p.l.Error("Deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJson(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJson(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}
