package repository

import (
    "database/sql"
    "github.com/MarcOrfilaCarreras/skyfare/internal/model"
)

type FlightsRepository struct {
    db *sql.DB
}

func NewFlightsRepository(db *sql.DB) *FlightsRepository {
    return &FlightsRepository{db: db}
}

func (r *FlightsRepository) InsertFlight(origin string, destination string, flight *model.Flight) error {
    _, err := r.db.Exec(`
		INSERT INTO flights (
			origin, destination, departure_date, price, promotion
		) VALUES (
			(SELECT id FROM airports WHERE code = ?),
            (SELECT id FROM airports WHERE code = ?),
			?, ?, ?)
		ON CONFLICT(origin, destination, departure_date) DO UPDATE SET
			price = excluded.price,
			departure_date = excluded.departure_date,
			promotion = excluded.promotion,
			updated_at = CURRENT_TIMESTAMP
		`,
        origin, destination, flight.Date, flight.Price, flight.Promotion,
    )
    return err
}
