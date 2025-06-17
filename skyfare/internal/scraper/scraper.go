package scraper

import (
    "github.com/MarcOrfilaCarreras/skyfare/internal/model"
)

type Scraper interface {
	GetAirports() ([]model.Airport, error)

	GetAirportRoutes(originCode string) ([]model.Route, error)

	GetRoute(
		originCode string,
		destCode string,
		month int,
		year int,
		currency string,
	) ([]model.Flight, error)
}