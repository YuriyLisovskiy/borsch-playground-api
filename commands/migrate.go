package commands

import (
	"github.com/YuriyLisovskiy/borsch-playground-service/migrations"
	"github.com/YuriyLisovskiy/borsch-playground-service/settings"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := settings.Load()
		if err != nil {
			return err
		}

		db, err := s.Database.Create()
		if err != nil {
			return err
		}

		return migrations.Migrate(db)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
