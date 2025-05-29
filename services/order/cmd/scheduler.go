package cmd

import (
	"context"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/order/internal"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/go-co-op/gocron/v2"
)

func startScheduler(logger *pkg.Logger, modules *internal.Modules) {
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Fatal("Failed to create scheduler:", err)
	}

	_, err = s.NewJob(gocron.CronJob("* * * * *", false), gocron.NewTask(func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelCtx()

		logger.Info("Processing expired orders...")
		if err := modules.Order.OrderService.ProcessExpiredOrders(ctx); err != nil {
			logger.Error("Failed to process expired orders:", err)
		}
	}))
	if err != nil {
		logger.Fatal("Failed to create job:", err)
	}

	logger.Info("Scheduler started, processing expired orders every minute")
	s.Start()
}
