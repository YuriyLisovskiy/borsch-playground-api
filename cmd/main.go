/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package main

import (
	"log"
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
	Use:  "borsch-playground-api",
	RunE: runRoot,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(
		&addressArg, "bind", "b", "127.0.0.1:8080", "bind address",
	)
}

func runRoot(*cobra.Command, []string) error {
	s, err := server.LoadSettingsFromEnv()
	if err != nil {
		return err
	}

	database, err := db.PostgreSQLFromEnv()
	if err != nil {
		return err
	}

	jobRepository := jobs.NewJobRepositoryImpl(database)
	amqpJobService := amqp.RabbitMQJobService{
		Server:        os.Getenv(amqp.EnvRabbitMQServer),
		JobRepository: jobRepository,
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

	jobService := jobs.NewJobServiceImpl(jobRepository, &amqpJobService)
	a := server.NewApplication(s, database, jobRepository, jobService)
	return a.Serve(addressArg)
}
