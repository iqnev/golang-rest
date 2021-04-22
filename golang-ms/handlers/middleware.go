package handlers

import (
	"context"
	"net/http"

	"github.com/iqnev/golang-rest/ms/data"
)

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		product := data.Product{}

		err := data.FromJson(product, r.Body)

		if err != nil {
			p.l.Error("Deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJson(&GenericError{Message: err.Error()}, rw)
			return
		}

		errs := p.v.Validate(product)

		if len(errs) != 0 {
			p.l.Error("Validating product", "error", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJson(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
