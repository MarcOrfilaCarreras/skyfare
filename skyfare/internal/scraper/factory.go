package scraper

import (
    "fmt"
    "strings"
    "database/sql"

    "github.com/MarcOrfilaCarreras/skyfare/internal/scraper/vueling"
)

func GetScraper(company string, db *sql.DB) (Scraper, error) {
    company = strings.ToLower(company)

    switch company {
    case "vueling":
        return vueling.NewVuelingScraper(db), nil
    default:
        return nil, fmt.Errorf("scraper for company '%s' not implemented", company)
    }
}