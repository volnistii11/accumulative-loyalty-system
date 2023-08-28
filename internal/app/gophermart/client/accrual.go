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

	logger.Info("access to the accumulation system is started")
	for {
		accrualService := service.NewAccrual()
		newOrders := accrualService.GetNewOrders(storage)
		if len(newOrders) > 0 {
			for _, newOrder := range newOrders {
				accrualSystemAddress := fmt.Sprintf("%s%s", cfg.GetAccrualSystemAddress(), "/api/orders/")
				answer, err := accrualService.SendOrderNumbersToAccrualSystem(newOrder, accrualSystemAddress)
				if err != nil {
					slog.Error("", sl.Err(err))
					continue
				}
				err = accrualService.UpdateAccrualInfoForOrderNumber(storage, answer)
				if err != nil {
					slog.Error("", sl.Err(err))
					continue
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}