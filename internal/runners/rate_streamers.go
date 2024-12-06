package runners

import (
	"context"
	"fmt"
	//sq "github.com/Masterminds/squirrel"
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

	_, err = conn.LoadCoins()
	if err != nil {
		return nil, fmt.Errorf("can't load coins")
	}

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

	rate, ok := data[""]
	if !ok {
		rate, ok = data["USDT"]
		if !ok {
			rate, ok = data["BTC"]
			if !ok {
				rate, ok = data["ETH"]
				if !ok {
					rate, ok = data["BNB"]
					if !ok {
						return fmt.Errorf("invalid rate data format: got %T", data["rate"])
					}
				}
			}
		}
	}

	//ON CONFLICT

	// SQL REQUEST
	query := "INSERT INTO rates (from_currency, to_currency, rate, provider, created_at, updated_at)" +
		" VALUES ($1, $2, $3, $4, $5, $6)" + " ON CONFLICT (from_currency, to_currency, provider)" +
		" DO UPDATE SET rate = EXCLUDED.rate, updated_at = EXCLUDED.updated_at;"
	_, err := r.db.Pool.Exec(context.Background(), query, r.rateFrom, r.rateTo, rate, r.connectorType, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}

	log.Println("Data saved to DB:", data)
	return nil

	/*	query, args, _ := sq.Insert("rates").
				Columns("from_currency", "to_currency", "rate", "provider", "created_at", "updated_at").
				Values(r.rateFrom, r.rateTo, rate, r.connectorType, time.Now(), time.Now()).
				Suffix("ON CONFLICT (from_currency, to_currency, provider) DO UPDATE SET rate = EXCLUDED.rate, updated_at = EXCLUDED.updated_at").
				ToSql()

			fmt.Println("SQL Query:", query)
			fmt.Println("Arguments:", args)

			conn, err := r.db.Pool.Acquire(r.db.Ctx)
			if err != nil {
				return fmt.Errorf("failed to acquire connection: %v", err)
			}
			defer conn.Release()

			_, err = conn.Exec(r.db.Ctx, query, args...)
			if err != nil {
				return fmt.Errorf("failed to insert data: %v", err)
			}

			log.Println("Data saved to DB:", data)
			return nil
		}
	*/
}
