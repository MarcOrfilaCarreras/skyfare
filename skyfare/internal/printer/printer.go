package printer

import (
    "encoding/json"

    "github.com/MarcOrfilaCarreras/skyfare/internal/model"
    "github.com/MarcOrfilaCarreras/skyfare/internal/logging"
)

func PrintAirports(jsonBytes []byte) error {
    var airports []model.Airport
    if err := json.Unmarshal(jsonBytes, &airports); err != nil {
        logging.Fatalf("failed to unmarshal airport JSON for printing: %w", err)
    }

    logging.Println("Available airports:")
    for _, airport := range airports {
        logging.Printf(" - %s (%s), Country: %s, Coordinates: %s, %s\n",
            airport.Name, airport.Code, airport.Country,
            airport.Coordinates.Lat, airport.Coordinates.Lng)
    }
    return nil
}

func PrintRoutes(jsonBytes []byte) error {
    var routes []model.Route
    if err := json.Unmarshal(jsonBytes, &routes); err != nil {
        logging.Fatalf("failed to unmarshal route JSON for printing: %w", err)
    }

	logging.Println("Available routes:")
    for _, route := range routes {
		logging.Printf(" - To: %s (Connection: %v)\n", route.Code, route.Connection)
	}
    return nil
}

func PrintFlights(jsonBytes []byte) error  {
    var flights []model.Flight
    if err := json.Unmarshal(jsonBytes, &flights); err != nil {
        logging.Fatalf("failed to unmarshal flight JSON for printing: %w", err)
    }

	logging.Println("Available flights:")
	for _, flight := range flights {
		promo := ""
		if flight.Promotion {
			promo = " (PROMO)"
		}
		logging.Printf(" - Date: %s | Price: %.2f%s\n", flight.Date, flight.Price, promo)
	}
    return nil
}