package cmd

import (
    "github.com/spf13/cobra"
    "github.com/MarcOrfilaCarreras/skyfare/internal/db"
    "github.com/MarcOrfilaCarreras/skyfare/internal/logging"
)

var (
    quiet bool
)

var rootCmd = &cobra.Command{
    Use:   "skyfare",
    CompletionOptions: cobra.CompletionOptions{
        DisableDefaultCmd: true,
    },
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        logging.SetQuiet(quiet)
        runMigrations()
    },
}

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}

func init() {
    rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output")
}

func runMigrations() {
    m, err := db.NewDatabaseMigrator("cache.db")
    if err != nil {
        logging.Fatalf("Failed to open database: %v", err)
    }
    defer m.Close()

    if err := m.Migrate(); err != nil {
        logging.Fatalf("Migration failed: %v", err)
    }
    logging.Println("Database migrations applied successfully.")
}
