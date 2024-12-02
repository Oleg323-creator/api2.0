package connectors

import (
	"fmt"
	"github.com/Oleg323-creator/api2.0/pkg/connectros/coingecko"
	"github.com/Oleg323-creator/api2.0/pkg/connectros/crypto_compare"
)

const СoingeckoType = "Coingecko"
const CryptoCompType = "Crypto Compare"

type ConnectorAPI interface {
	LoadCoins() (int, error)
	GetRates(from, to string) (map[string]interface{}, error)
}

func NewConnector(conType string) (ConnectorAPI, error) {

	if conType == СoingeckoType {
		return coingecko.NewGeckoApi(), nil
	} else if conType == CryptoCompType {
		return crypto_compare.NewCryptoCompareAPI(), nil
	} else {
		return nil, fmt.Errorf("unknown connector type")
	}
}

/*func main() {
	//GECKO IMPLEMENTATION
	conn, err := NewConnector(СoingeckoType)
	if err != nil {
		return
	}
	_, err = conn.LoadCoins()
	if err != nil {
		return
	}

	rate, err := conn.GetRates("BTC", "USDT")
	if err != nil {
		return
	}
	fmt.Println(rate)

	//CRYPTO_COMPARE IMPLEMENTATION
	conn, err = NewConnector(CryptoCompType)
	if err != nil {
		return
	}
	_, err = conn.LoadCoins()
	if err != nil {
		return
	}

	rate, err = conn.GetRates("BTC", "USDT")
	if err != nil {
		return
	}
	fmt.Println(rate)
}*/
