package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bleedsix/web-portal/internal/data" // Заміни на реальний шлях до пакету data
)

// tournamentsRepo — це приватна структура, яка реалізує інтерфейс data.TournamentsRepository
type tournamentsRepo struct {
	db *sql.DB
}

// NewTournamentsRepo створює новий екземпляр репозиторію
func NewTournamentsRepo(db *sql.DB) *tournamentsRepo {
	return &tournamentsRepo{db: db}
}

// Get повертає турнір за ID
func (r *tournamentsRepo) Get(ctx context.Context, id string) (*data.Tournament, error) {
	query := `
		SELECT id, name, format, lan, region, prize_pool, prize_distribution, date
		FROM tournaments
		WHERE id = $1
	`

	t := &data.Tournament{}
	// Скануємо дані. Зверни увагу: prize_distribution скануємо як []byte, бо це JSONB
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Name,
		&t.Format,
		&t.Lan,
		&t.Region,
		&t.PrizePool,
		&t.PrizeDistribution,
		&t.Date,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tournament not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get tournament: %w", err)
	}

	return t, nil
}

// Upsert оновлює існуючий запис або вставляє новий (ON CONFLICT)
func (r *tournamentsRepo) Upsert(ctx context.Context, t *data.Tournament) error {
	query := `
		INSERT INTO tournaments (id, name, format, lan, region, prize_pool, prize_distribution, date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			format = EXCLUDED.format,
			lan = EXCLUDED.lan,
			region = EXCLUDED.region,
			prize_pool = EXCLUDED.prize_pool,
			prize_distribution = EXCLUDED.prize_distribution,
			date = EXCLUDED.date;
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.Name,
		t.Format,
		t.Lan,
		t.Region,
		t.PrizePool,
		t.PrizeDistribution, // []byte або json.RawMessage нормально лягає в JSONB
		t.Date,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert tournament: %w", err)
	}

	return nil
}