package runners

import (
	"context"
	"fmt"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/Oleg323-creator/api2.0/pkg/connectros"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	db            *db.WrapperDB
}

func NewRunner(conType string, pollRate int, from string, to string, db *db.WrapperDB) (*Runner, error) {
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
		db:            db,
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
				log.Printf("Error fetching rates: %v", err)
				continue
			}

			//SAIVING RATES TO DB
			err = r.saveDataToDB(rates)
			if err != nil {
				log.Printf("Error saving data to DB: %v", err)
			}

			log.Println(time.Now().Unix(), r.rateFrom, rates, r.connectorType)
			continue

		}
	}
}

func (r *Runner) saveDataToDB(data map[string]interface{}) error {

	rate, ok := data["USDT"].(float64)
	if !ok {
		return fmt.Errorf("invalid rate data format: expected float64, got %T", data["rate"])
	}

	// SQL REQUEST
	query := "INSERT INTO rates (from_currency, to_currency, rate, provider, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Pool.Exec(context.Background(), query, r.rateFrom, r.rateTo, rate, r.connectorType, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}

	log.Println("Data saved to DB:", data)
	return nil
}
