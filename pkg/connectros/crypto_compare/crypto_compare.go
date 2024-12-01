package crypto_compare

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var ticker = map[string]string{
	"BTC":  "bitcoin",
	"ETH":  "ethereum",
	"USDT": "usd",
}

// LOADING .ENV
func init() {
	err := godotenv.Load("/home/oleh/GolandProjects/geckoapi1/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

type CryptoCompareAPI struct {
	URL string
	KEY string
}

func NewCryptoCompareAPI() *CryptoCompareAPI {
	return &CryptoCompareAPI{
		URL: "https://min-api.cryptocompare.com/data",
		KEY: os.Getenv("API_KEY"),
	}
}

func (c *CryptoCompareAPI) LoadCoins() (int, error) {
	return 0, nil
}

func (c *CryptoCompareAPI) GetRates(from, to string) (map[string]interface{}, error) {
	client := &http.Client{}
	endpoint := "/price"

	//MAKING REQUEST URL
	reqURL := c.URL + endpoint
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request:%v", err)
	}

	//ADDING PARAMETRS TO URL
	params := map[string]string{
		"fsym":  from,
		"tsyms": to,
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	//SENDING REQUEST
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	//READING BODY RESPONSE
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	//fmt.Println("Response Body:", string(body))

	//CHECKING STATUS CODE
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	//PARSING JSON
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %v", err)
	}

	//CHEKICNG IF ERROR
	if response, ok := result["Response"]; ok && response == "Error" {
		if message, exists := result["Message"].(string); exists {
			return nil, fmt.Errorf("API error: %s", message)
		}
		return nil, fmt.Errorf("API error occurred without message")
	}

	return result, nil
}
