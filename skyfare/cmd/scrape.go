package cmd

import (
    "encoding/json"
    "time"
    "database/sql"

    "github.com/spf13/cobra"
    "github.com/MarcOrfilaCarreras/skyfare/internal/scraper"
    "github.com/MarcOrfilaCarreras/skyfare/internal/printer"
    "github.com/MarcOrfilaCarreras/skyfare/internal/logging"
)

var (
    company string
    command string
    origin string
    destination string
    month int
    year int
    currency string
)

var allowedCompanies = []string{"vueling"}
var allowedCommands = []string{"airports", "routes", "route"}

var scrapeCmd = &cobra.Command{
    Use:   "scrape",
    Short: "Run the scraper",
    Run: func(cmd *cobra.Command, args []string) {
        if company == "" || command == "" {
            logging.Fatalf("All parameters are required: --company, --command")
        }

        db, err := sql.Open("sqlite3", "cache.db")
        if err != nil {
            logging.Fatalf("Error opening database: %v", err)
        }
        defer db.Close()

        s, err := scraper.GetScraper(company, db)
        if err != nil {
            logging.Fatalf("Error initializing scraper: %v", err)
        }

        switch command {
        case "airports":
            airports, err := s.GetAirports()
            if err != nil {
                logging.Fatalf("Error fetching airports: %v", err)
            }

            output, err := json.MarshalIndent(airports, "", "  ")
            if err != nil {
                logging.Fatalf("Error formatting output: %v", err)
            }

            printer.PrintAirports(output)

        case "routes":
            if origin == "" {
                logging.Fatalf("--origin is required for routes command")
            }

            routes, err := s.GetAirportRoutes(origin)
            if err != nil {
                logging.Fatalf("Error fetching routes: %v", err)
            }

            output, err := json.MarshalIndent(routes, "", "  ")
            if err != nil {
                logging.Fatalf("Error formatting output: %v", err)
            }

            printer.PrintRoutes(output)

        case "flights":
            if origin == "" {
                logging.Fatalf("--origin is required for route command")
            }

            if destination == "" {
                logging.Fatalf("--destination is required for route command")
            }

            flights, err := s.GetRoute(origin, destination, month, year, currency)
            if err != nil {
                logging.Fatalf("Error fetching flights: %v", err)
            }

            output, err := json.MarshalIndent(flights, "", "  ")
            if err != nil {
                logging.Fatalf("Error formatting output: %v", err)
            }

            printer.PrintFlights(output)

        default:
            logging.Fatalf("Unknown command: %s", command)
        }
    },
}

func init() {
    now := time.Now()
    month = int(now.Month())
    year = now.Year()

    rootCmd.AddCommand(scrapeCmd)
    scrapeCmd.Flags().StringVar(&company, "company", "", "Flight company name (required)")
    scrapeCmd.Flags().StringVar(&command, "command", "", "Command to execute (airports, routes, flights)")
    scrapeCmd.Flags().StringVar(&origin, "origin", "", "Origin airport code (required for routes command)")
    scrapeCmd.Flags().StringVar(&destination, "destination", "", "Destination airport code (required for routes command)")
    scrapeCmd.Flags().IntVar(&month, "month", month, "Month (required for routes command)")
    scrapeCmd.Flags().IntVar(&year, "year", year, "Year (required for routes command)")
    scrapeCmd.Flags().StringVar(&currency, "currency", "EUR", "Currency code (required for routes command)")

    scrapeCmd.MarkFlagRequired("company")
    scrapeCmd.MarkFlagRequired("command")
}
