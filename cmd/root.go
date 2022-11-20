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

	"borsch-playground-api/app"
	"borsch-playground-api/jobs"
	rmq "borsch-playground-api/rmq"
	"borsch-playground-api/settings"
	"github.com/spf13/cobra"
)

var (
	addressArg string
)

var rootCmd = &cobra.Command{
	Use:  "borsch-playground-api",
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
	s, err := settings.Load()
	if err != nil {
		return err
	}

	db, err := s.Database.Build()
	if err != nil {
		return err
	}

	jobService := jobs.NewJobServiceImpl(db)
	amqpJobService := rmq.RabbitMQJobService{
		Server:     os.Getenv(rmq.EnvRabbitMQServer),
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

	a, err := app.NewApp(s, db, jobService, &amqpJobService)
	if err != nil {
		return err
	}

	return a.Execute(addressArg)
}
