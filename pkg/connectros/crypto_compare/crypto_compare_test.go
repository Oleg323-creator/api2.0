package crypto_compare

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRates(t *testing.T) {
	api := NewCryptoCompareAPI()

	// GETTING RATE FOR BTC

	data, err := api.GetRates("BTC", "USDT")

	assert.NoError(t, err, "Unexpected error occurred")

	if assert.NotNil(t, data, "Expected data but got nil") {
		assert.Contains(t, data, "USDT", "Expected USD in response")
	}

	fmt.Println(data)

	// GETTING RATE FOR ETH

	data, err = api.GetRates("ETH", "USDT")

	assert.NoError(t, err, "Unexpected error occurred")

	if assert.NotNil(t, data, "Expected data but got nil") {
		assert.Contains(t, data, "USDT", "Expected USDT in response")
	}

	fmt.Println(data)

}

// TEST FOR INCORRECT CURR
func TestGetRatesError(t *testing.T) {
	api := NewCryptoCompareAPI()

	// RANDOM RATE

	data, err := api.GetRates("nonexistentcoin", "USD")

	// EXPECTED ERROR
	if assert.Error(t, err, "Expected an error but got nil") {
		assert.Contains(t, err.Error(), "market does not exist for this coin pair", "Error message mismatch")
	}

	assert.Nil(t, data, "Expected data to be nil on error")
}
