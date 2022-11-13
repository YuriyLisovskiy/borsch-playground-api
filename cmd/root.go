/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package cmd

import (
	"os"

	"github.com/YuriyLisovskiy/borsch-playground-api/app"
	"github.com/YuriyLisovskiy/borsch-playground-api/jobs"
	rmq "github.com/YuriyLisovskiy/borsch-playground-api/rabbitmq"
	"github.com/YuriyLisovskiy/borsch-playground-api/settings"
	"github.com/spf13/cobra"
)

var (
	addressArg string
)

var rootCmd = &cobra.Command{
	Use: "borsch-playground-api",
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := settings.Load()
		if err != nil {
			return err
		}

		db, err := settings.Database.Create()
		if err != nil {
			return err
		}

		jobService := jobs.NewJobServiceImpl(db)
		amqpJobService := rmq.RabbitMQJobService{
			Server:     os.Getenv("RABBITMQ_SERVER"),
			JobService: jobService,
		}
		err = amqpJobService.Setup()
		if err != nil {
			return err
		}

		defer amqpJobService.CleanUp()
		err = amqpJobService.ConsumeJobResults()
		if err != nil {
			return err
		}

		app, err := app.NewApp(settings, db, jobService, &amqpJobService)
		if err != nil {
			return err
		}

		return app.Execute(addressArg)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(
		&addressArg, "address", "a", "127.0.0.1:8080", "server address",
	)
}
