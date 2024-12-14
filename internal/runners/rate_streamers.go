package runners

import (
	"context"
	"fmt"
	"github.com/Oleg323-creator/api2.0/internal/db/rep"
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
	repository    *rep.Repository
}

func NewRunner(conType string, pollRate int, from string, to string, repository *rep.Repository) (*Runner, error) {
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
		repository:    repository,
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
			err = r.repository.SaveDataToDB(r.rateFrom, r.rateTo, r.connectorType, rates)
			if err != nil {
				log.Printf("Error saving data to DB: %v", err)
			}

			log.Println(time.Now().Unix(), r.rateFrom, rates, r.connectorType)
			continue

		}
	}
}
