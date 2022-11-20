/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package cmd

import (
	"borsch-playground-api/migrations"
	"borsch-playground-api/settings"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	RunE:  migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate(*cobra.Command, []string) error {
	s, err := settings.Load()
	if err != nil {
		return err
	}

	db, err := s.Database.Build()
	if err != nil {
		return err
	}

	return migrations.Migrate(db)
}
