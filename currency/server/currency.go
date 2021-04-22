package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iqnev/golang-rest/currency/data"
	protos "github.com/iqnev/golang-rest/currency/protos/currency"
)

type Currency struct {
	log           hclog.Logger
	retes         *data.ExchangeRates
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(l hclog.Logger, r *data.ExchangeRates) *Currency {
	c := &Currency{l, r, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}

	go c.handleUpdates()

	return c
}

func (c *Currency) handleUpdates() {
	ru := c.retes.MonitorRates(3 * time.Second)

	for range ru {
		c.log.Info("Got Updated rates")

		for k, v := range c.subscriptions {
			for _, rr := range v {
				r, err := c.retes.GetRate(rr.GetBase().String(), rr.GetDestination().String())

				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{RateResponse: &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r}},
				})
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, in *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", in.GetBase(), "dest", in.GetDestination())

	if in.Base == in.Destination {
		err := status.Errorf(
			codes.InvalidArgument,
			"Base rate %s can not be equal to destination rate %s",
			in.Base.String(),
			in.Destination.String(),
		)

		return nil, err
	}
	rate, err := c.retes.GetRate(in.GetBase().String(), in.GetDestination().String())

	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{
		Base:        in.Base,
		Destination: in.Destination,
		Rate:        rate,
	}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	for {
		rcv, errcv := src.Recv()

		if errcv == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}

		if errcv != nil {
			c.log.Error("Unable to read from client", "error", errcv)
			break
		}

		c.log.Info("Handle client request", "request_base", rcv.GetBase(), "request_dest", rcv.GetDestination())

		rrs, ok := c.subscriptions[src]

		if !ok {
			rrs = []*protos.RateRequest{}
		}

		for _, r := range rrs {
			if r.Base == rcv.Base && r.Destination == rcv.Destination {
				c.log.Error("Subscription already active", "base", rcv.Base.String(), "dest", rcv.Destination.String())

				grpcError := status.New(codes.InvalidArgument, "Subscription already active for rate")
				grpcError, errcv = grpcError.WithDetails(rcv)

				if errcv != nil {
					c.log.Error("Unable to add metadata to error message", "error", errcv)
					continue
				}

				rrs := protos.StreamingRateResponse_Error{Error: grpcError.Proto()}

				src.Send(&protos.StreamingRateResponse{Message: &rrs})

			}
		}

		rrs = append(rrs, rcv)
		c.subscriptions[src] = rrs
	}

	return nil
}
