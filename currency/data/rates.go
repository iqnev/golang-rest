package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

func NewExchangeRates(l hclog.Logger) (*ExchangeRates, error) {
	exr := &ExchangeRates{log: l, rates: map[string]float64{}}

	err := exr.getRates()
	return exr, err
}

func (exr *ExchangeRates) GetRate(base string, dest string) (float64, error) {
	br, ok := exr.rates[base]

	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}

	dr, ok := exr.rates[dest]

	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}

	return dr / br, nil
}

func (exr *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	newCubes := &Cubes{}

	xml.NewDecoder(resp.Body).Decode(newCubes)

	for _, c := range newCubes.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		exr.rates[c.Currency] = r
	}

	exr.rates["EUR"] = 1

	return nil
}

func (exr *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)

		for {
			select {
			// just add a random difference to the rate and return it
			// this simulates the fluctuations in currency rates
			case <-ticker.C:
				for k, v := range exr.rates {
					change := (rand.Float64() / 10)

					direction := rand.Intn(1)

					if direction == 0 {
						change = 1 - change
					} else {
						change = 1 + change
					}

					// modify the rate
					exr.rates[k] = v * change
				}

				ret <- struct{}{}
			}
		}
	}()

	return ret
}
