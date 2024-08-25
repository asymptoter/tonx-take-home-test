package main

import (
	"github.com/asymptoter/tonx-take-home-test/internal/api/handler"
	"github.com/asymptoter/tonx-take-home-test/internal/repository"
	"github.com/asymptoter/tonx-take-home-test/internal/service"
	"github.com/asymptoter/tonx-take-home-test/pkg/ctx"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	ctx := ctx.Background()
	// Connect to database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		ctx.Fatal(err)
	}

	campaignRepository := repository.NewCampaignRepository(ctx, db)
	campaignService := service.NewCampaignService(ctx, campaignRepository)

	// Cron job create campaign every day
	cronJob := cron.New(cron.WithSeconds())
	if _, err = cronJob.AddFunc("0 30 22 * * *", func() {
		if _, err := campaignService.Create(ctx, service.CreateCampaignInput{}); err != nil {
			ctx.Fatal(err)
		}
	}); err != nil {
		ctx.Fatal(err)
	}
	cronJob.Start()

	router := gin.Default()
	handler.RegisterHTTPHandler(router, campaignService)
	router.Run(":8080")
}
