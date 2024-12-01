package coingecko

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadCoins(t *testing.T) {
	api := NewGeckoApi()

	// AMOUNT OF COINS
	count, err := api.LoadCoins()

	// CHECKING ERROR
	assert.NoError(t, err)

	// CHECKING IF AMOUNT > 0
	assert.True(t, count > 0, "Expected to have available coins")
}

func TestGetRates(t *testing.T) {

	api := NewGeckoApi()

	// GETTING RATE FOR BTC
	data, err := api.GetRates("BTC", "USDT")

	assert.NoError(t, err)

	fmt.Println(data)

	// GETTING RATE FOR ETH
	data, err = api.GetRates("ETH", "USDT")

	assert.NoError(t, err)

	fmt.Println(data)

}

// TEST FOR INCORRECT CURR
func TestGetRatesError(t *testing.T) {
	api := NewGeckoApi()

	// RANDOM RATE
	_, err := api.GetRates("nonexistentcoin", "USDT")

	// EXPECTED ERROR
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "conversion rate not found")
}
