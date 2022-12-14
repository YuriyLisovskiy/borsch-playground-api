/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"borsch-playground-api/jobs"
	rmq "borsch-playground-api/rmq"
	"borsch-playground-api/settings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	settings       *settings.Settings
	db             *gorm.DB
	jobService     jobs.JobService
	amqpJobService rmq.AMQPJobService
}

func NewApp(
	s *settings.Settings,
	db *gorm.DB,
	jobService jobs.JobService,
	amqpJobService rmq.AMQPJobService,
) (*Application, error) {
	gin.SetMode(s.GinMode)
	app := &Application{
		settings:       s,
		db:             db,
		jobService:     jobService,
		amqpJobService: amqpJobService,
	}
	return app, nil
}

func (a *Application) buildRouter() *gin.Engine {
	router := gin.Default()
	a.addV1Routes(router)
	return router
}

func (a *Application) Execute(addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := a.buildRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), a.settings.ShutdownTimeoutSec*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return errors.New(fmt.Sprintf("Server forced to shut down: %v", err))
	}

	log.Println("Server exiting")
	return nil
}
