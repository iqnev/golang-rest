package handlers

import (
	"net/http"

	"github.com/iqnev/golang-rest/ms/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, req *http.Request) {
	p.l.Debug("Get all records")
	rw.Header().Add("Content-Type", "application/json")

	cr := req.URL.Query().Get("currency")
	products, err := p.productDB.GetProducts(cr)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJson(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJson(products, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serializing product", "error", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse
func (p *Products) ListSingle(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(req)

	cur := req.URL.Query().Get("currency")

	p.l.Debug("Get record", "id", id)

	prod, err := p.productDB.GetProductByID(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Unable to fetch product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJson(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Unable to fetching product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJson(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJson(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to  serializing product", err)
	}
}
