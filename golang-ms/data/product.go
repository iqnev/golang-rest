package data

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	protos "github.com/iqnev/golang-rest/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrProductNotFound = fmt.Errorf("Product not found")

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"  validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"  validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeleteOn    string  `json:"-"`
}

type Products []*Product

type ProductDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductDB(c protos.CurrencyClient, l hclog.Logger) *ProductDB {
	pb := &ProductDB{c, l, map[string]float64{}, nil}

	go pb.handleUpdates()

	return pb
}

func (pr *ProductDB) handleUpdates() {
	sub, err := pr.currency.SubscribeRates(context.Background())

	if err != nil {
		pr.log.Error("Unable to subscribe for rates", "error", err)
	}

	pr.client = sub

	for {
		rr, err := sub.Recv()
		pr.log.Info("Recieved updated rate from server", "dest", rr.GetDestination().String())

		if err != nil {
			pr.log.Error("Error receiving message", "error", err)
			return
		}

		pr.rates[rr.Destination.String()] = rr.Rate
	}
}

func (pr *ProductDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := pr.getRate(currency)

	if err != nil {
		pr.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	prd := Products{}

	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		prd = append(prd, &np)

	}

	return prd, nil
}

func (pr *ProductDB) getRate(destination string) (float64, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies_EUR,
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := pr.currency.GetRate(context.Background(), rr)

	if err != nil {
		grpsErr, ok := status.FromError(err)

		if !ok {
			return -1. err
		}

		if grpsErr.Code() == codes.InvalidArgument {
			return -1, fmt.Errorf("Unable to retreive exchange rate from currency service: %s", grpcError.Message())
		}

	}

	//update cache
	pr.rates[destination] = resp.Rate

	pr.client.Send(rr)

	return resp.Rate, err
}

func (p *ProductDB) UpdateProduct(pr *Product) error {
	i := findIndexByProductID(pr.ID)

	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = pr

	return nil
}

func (p *ProductDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)

	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

func findIndexByProductID(id int) int {

	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(currency)

	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i]

	np.Price = np.Price * rate

	return &np, nil
}

func (p *ProductDB) AddProduct(pr *Product) {
	pr.ID = getNextID()

	productList = append(productList, pr)
}

func getNextID() int {
	lp := productList[len(productList)-1]

	return lp.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Strong milk coffee",
		Price:       1.99,
		SKU:         "def123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
