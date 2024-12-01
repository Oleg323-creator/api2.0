package main

import (
	"context"
	"fmt"
	"geckoapi1/pkg/connectors"
	"log"
	"sync"
	"time"
)

type Runner struct {
	connectorInit connectors.ConnectorAPI
	connectorType string
	pollingRate   int
	rateFrom      string
	rateTo        string
}

func NewRunner(conType string, pollRate int, from string, to string) (*Runner, error) {
	conn, err := connectors.NewConnector(conType)
	if err != nil {
		return nil, fmt.Errorf("invalid connector type")
	}
	coins, err := conn.LoadCoins()
	if err != nil {
		return nil, fmt.Errorf("can't load coins")
	}
	fmt.Println(coins)

	return &Runner{
		connectorInit: conn,
		connectorType: conType,
		pollingRate:   pollRate,
		rateFrom:      from,
		rateTo:        to,
	}, nil
}

func (r *Runner) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(r.pollingRate) * time.Second)
	log.Println("Starting:")
	for {
		select {
		case <-ctx.Done():
			log.Println("Finishing:")
			return
		case <-ticker.C:
			rates, err := r.connectorInit.GetRates(r.rateFrom, r.rateTo)
			if err != nil {
				return
			}

			log.Println(time.Now().Unix(), r.rateFrom, rates, r.connectorType)
			continue
		}
	}
}

/*func main() {
	geckoRun, err := NewRunner(connectors.СoingeckoType, 3, "BTC", "USDT")
	if err != nil {
		return
	}
	cryptCompRun, err := NewRunner(connectors.CryptoCompType, 1, "ETH", "USDT")
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go geckoRun.Run(ctx, wg)
	go cryptCompRun.Run(ctx, wg)

	<-stop

	cancel()
	wg.Wait()
}*/
