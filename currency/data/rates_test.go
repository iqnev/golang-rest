package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewrates(t *testing.T) {

	tr, err := NewExchangeRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Rates %#v", tr.rates)
}
