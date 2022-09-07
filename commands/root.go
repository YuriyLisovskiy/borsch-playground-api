/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package commands

import (
	"github.com/YuriyLisovskiy/borsch-playground-service/app"
	"github.com/YuriyLisovskiy/borsch-playground-service/settings"
	"github.com/spf13/cobra"
)

var (
	addressArg string
)

var rootCmd = &cobra.Command{
	Use: "borschplayground",
	RunE: func(cmd *cobra.Command, args []string) error {
		appSettings, err := settings.Load()
		if err != nil {
			return err
		}

		app, err := app.NewApp(appSettings)
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
