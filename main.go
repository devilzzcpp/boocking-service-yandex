package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"booking_service/config"
	"booking_service/internal/app"
	"booking_service/internal/entity/booking"
	"booking_service/utils"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if port != 8080 {
		cfg.Port = port
	}

	lg, err := utils.New(cfg.AppEnv)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer func() {
		_ = lg.Sync()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := utils.Connect(ctx, cfg)
	if err != nil {
		lg.Fatal("connect db", utils.Error(err))
	}
	defer func() {
		if err := utils.Close(db); err != nil {
			lg.Error("close db", utils.Error(err))
		}
	}()

	repo := booking.NewRepository(db)
	svc := booking.NewService(repo)
	h := booking.NewHandler(svc)

	r := app.NewRouter(lg, h)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		lg.Fatal("run server", utils.Error(err))
	}
}
