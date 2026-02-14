package data

import (
	"context"
	"encoding/json"
	"time"
)

// Tournament відображає структуру таблиці tournaments
type Tournament struct {
	ID                string          `json:"id"`
	Name              *string         `json:"name"` // В SQL це nullable, тому вказівник
	Format            string          `json:"format"`
	Lan               bool            `json:"lan"`
	Region            string          `json:"region"`
	PrizePool         *float64        `json:"prize_pool"` // Numeric nullable
	PrizeDistribution json.RawMessage `json:"prize_distribution"` // JSONB
	Date              time.Time       `json:"date"`
}

// TournamentsRepository описує контракт (інтерфейс) для роботи з турнірами
type TournamentsRepository interface {
	Get(ctx context.Context, id string) (*Tournament, error)
	Upsert(ctx context.Context, t *Tournament) error
}