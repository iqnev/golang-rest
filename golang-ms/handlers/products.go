package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/iqnev/golang-rest/ms/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l         hclog.Logger
	v         *data.Validation
	productDB *data.ProductDB
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

type KeyProduct struct{}

var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductDB) *Products {
	return &Products{l, v, pdb}
}

func getProductID(r *http.Request) int {
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	return id
}
