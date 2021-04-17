package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/iqnev/golang-rest/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
	v *data.Validation
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

type KeyProduct struct{}

var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
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
