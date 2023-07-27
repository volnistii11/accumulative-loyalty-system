package accrual

import (
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
			for _, newOrder := range newOrders {
				answer, err := accrualService.SendOrderNumbersToAccrualSystem(newOrder, cfg.GetAccrualSystemAddress())
				if err != nil {
					slog.Error("", sl.Err(err))
					return
				}
				err = accrualService.UpdateAccrualInfoForOrderNumber(storage, answer)
				if err != nil {
					slog.Error("", sl.Err(err))
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
