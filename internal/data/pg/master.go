package pg

import (
	"database/sql"
	
	"github.com/bleedsix/web-portal/internal/data"
)

// Storage — це реалізація інтерфейсу data.Storage для PostgreSQL
type Storage struct {
	db          *sql.DB
	tournaments data.TournamentsRepository
}

// NewStorage створює нове сховище з підключенням до БД
func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db:          db,
		tournaments: NewTournamentsRepo(db), // Ініціалізуємо конкретну реалізацію
	}
}

// Tournaments повертає інтерфейс для роботи з турнірами
func (s *Storage) Tournaments() data.TournamentsRepository {
	return s.tournaments
}