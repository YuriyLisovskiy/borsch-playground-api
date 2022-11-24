/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package cli

import (
	"os"

	"github.com/YuriyLisovskiy/borsch-playground-api/internal/amqp"
	"github.com/YuriyLisovskiy/borsch-playground-api/internal/db"
	"github.com/YuriyLisovskiy/borsch-playground-api/internal/jobs"
	"github.com/YuriyLisovskiy/borsch-playground-api/internal/server"
	"github.com/spf13/cobra"
)

var (
	addressArg string
)

var rootCmd = &cobra.Command{
	Use:  "github.com/YuriyLisovskiy/borsch-playground-api",
	RunE: root,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(
		&addressArg, "bind", "b", "127.0.0.1:8080", "bind address",
	)
}

func root(*cobra.Command, []string) error {
	s, err := server.LoadSettingsFromEnv()
	if err != nil {
		return err
	}

	database, err := db.PostgreSQLFromEnv()
	if err != nil {
		return err
	}

	jobService := jobs.NewJobServiceImpl(database)
	amqpJobService := amqp.RabbitMQJobService{
		Server:     os.Getenv(amqp.EnvRabbitMQServer),
		JobService: jobService,
	}
	err = amqpJobService.Setup()
	if err != nil {
		return err
	}

	defer amqpJobService.CleanUp()
	err = amqpJobService.ConsumeResults()
	if err != nil {
		return err
	}

	a := server.NewApplication(s, database, jobService, &amqpJobService)
	return a.Serve(addressArg)
}
