package accrual

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/volnistii11/accumulative-loyalty-system/internal/app/gophermart/service"
	"github.com/volnistii11/accumulative-loyalty-system/internal/config"
	"github.com/volnistii11/accumulative-loyalty-system/internal/lib/sl"
	"github.com/volnistii11/accumulative-loyalty-system/internal/repository/database"
	"golang.org/x/exp/slog"
	"net/http"
)

func DoAccrualIfPossible(logger *slog.Logger, storage *database.Storage, cfg config.ParserGetter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			const destination = "middleware.accrual.DoAccrualIfPossible"
			logger = logger.With(
				slog.String("destination", destination),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			accrualService := service.NewAccrual()
			newOrders := accrualService.GetNewOrders(storage)
			fmt.Println("------------------------------------------")
			fmt.Println("NEW ORDERS", newOrders)
			fmt.Println("LEN NEW ORDERS", len(newOrders))
			fmt.Println("ACCRUAL SYSTEM ADDRESS", cfg.GetAccrualSystemAddress())
			if len(newOrders) > 0 {
				for _, newOrder := range newOrders {
					answer, err := accrualService.SendOrderNumbersToAccrualSystem(newOrder, cfg.GetAccrualSystemAddress())
					fmt.Println("ANSWER", answer)
					fmt.Println("ANSWERERROR", err)
					if err != nil {
						slog.Error("", sl.Err(err))
						continue
					}
					err = accrualService.UpdateAccrualInfoForOrderNumber(storage, answer)
					fmt.Println("UPDATEERR", err)
					if err != nil {
						slog.Error("", sl.Err(err))
						continue
					}
				}
			}
			fmt.Println("------------------------------------------")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
