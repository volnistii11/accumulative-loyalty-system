package client

import (
	"fmt"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
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
			for _, newOrder := range newOrders {
				accrualSystemAddress := fmt.Sprintf("%s%s", cfg.GetAccrualSystemAddress(), "/api/orders/")
				answer, err := accrualService.SendOrderNumbersToAccrualSystem(newOrder, accrualSystemAddress)
				if err != nil {
					logger.Error("", sl.Err(err))
					continue
				}
				logger.Info("accrual system answer: ", answer)
				err = accrualService.UpdateAccrualInfoForOrderNumber(storage, answer)
				if err != nil {
					logger.Error("", sl.Err(err))
					continue
				}
			}
		}
		time.Sleep(time.Second)
	}
}
