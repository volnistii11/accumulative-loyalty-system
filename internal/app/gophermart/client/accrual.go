package client

import (
	"errors"
	"fmt"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/cerrors"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
	"sync"
	"time"
)

func DoAccrualIfPossible(logger *slog.Logger, storage *database.Storage, cfg config.ParserGetter) {
	const destination = "client.accrual.DoAccrualIfPossible"
	logger = logger.With(
		slog.String("destination", destination),
	)

	for {
		accrualService := service.NewAccrual()
		newOrders := accrualService.GetNewOrders(storage)
		if len(newOrders) > 0 {
			var wg sync.WaitGroup

			for _, newOrder := range newOrders {
				wg.Add(1)

				go func(newOrder string) {
					accrualSystemAddress := fmt.Sprintf("%s%s", cfg.GetAccrualSystemAddress(), "/api/orders/")
					answer, err := accrualService.SendOrderNumbersToAccrualSystem(newOrder, accrualSystemAddress)
					if err != nil {
						logger.Error("", sl.Err(err))
						if errors.Is(err, cerrors.ErrHTTPStatusTooManyRequests) {
							time.Sleep(time.Second * 60)
						}
					} else {
						logger.Info("accrual system answer: ", answer)
						err = accrualService.UpdateAccrualInfoForOrderNumber(storage, answer)
						if err != nil {
							logger.Error("", sl.Err(err))
						}
					}
				}(newOrder)

			}

			wg.Wait()
		}
		time.Sleep(time.Second * 10)
	}
}
