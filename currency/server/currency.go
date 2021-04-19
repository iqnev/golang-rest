package server

import (
	"context"

	"github.com/hashicorp/go-hclog"

	protos "github.com/iqnev/golang-rest/currency/protos/currency"
)

type Currency struct {
	log hclog.Logger
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{l}
}

func (c *Currency) GetRate(ctx context.Context, in *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", in.GetBase(), "dest", in.GetDestination())

	return &protos.RateResponse{Rate: 0.33}, nil
}
