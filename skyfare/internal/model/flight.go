package model

import (
	"time"
)

type Flight struct {
	Date      time.Time	`json:"date"`
	Price     float64	`json:"price,omitempty"`
	Promotion bool		`json:"promotion"`
}
