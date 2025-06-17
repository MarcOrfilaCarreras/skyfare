package repository

import (
    "database/sql"
    "github.com/MarcOrfilaCarreras/skyfare/internal/model"
)

type AirportRepository struct {
    db *sql.DB
}

func NewAirportRepository(db *sql.DB) *AirportRepository {
    return &AirportRepository{db: db}
}

func (r *AirportRepository) InsertAirport(a model.Airport) error {
    _, err := r.db.Exec(`
        INSERT OR IGNORE INTO airports (code, name, country, lat, lng)
        VALUES (?, ?, ?, ?, ?);
    `, a.Code, a.Name, a.Country, a.Coordinates.Lat, a.Coordinates.Lng)
    return err
}
