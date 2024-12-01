package coingecko

import (
	"fmt"
	"log"
	"strings"

	gecko "github.com/superoo7/go-gecko/v3"
)

type GeckoApi struct {
	URL string
}

func NewGeckoApi() *GeckoApi {
	return &GeckoApi{
		URL: "https://api.coingecko.com/api/v3",
	}
}

var cg = gecko.NewClient(nil)

func (g *GeckoApi) LoadCoins() (int, error) {
	list, err := cg.CoinsList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Available coins:", len(*list))

	return len(*list), nil
}

func (g *GeckoApi) GetRates(from, to string) (map[string]interface{}, error) {
	var ticker = map[string]string{
		"BTC":  "bitcoin",
		"ETH":  "ethereum",
		"USDT": "usd",
	}
	from = ticker[from]
	to = ticker[to]

	fromParam := []string{strings.ToLower(from)}
	toParam := []string{strings.ToLower(to)}
	result := make(map[string]interface{})

	sp, err := cg.SimplePrice(fromParam, toParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}

	if rate, ok := (*sp)[from][to]; ok {
		result[ticker[to]] = float64(rate)
		return result, nil
	} else {
		return nil, fmt.Errorf("conversion rate not found")
	}
}
