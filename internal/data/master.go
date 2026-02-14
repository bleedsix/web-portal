package data

// Storage (або Master) — це головний інтерфейс, який дає доступ до всіх під-репозиторіїв.
type Storage interface {
	Tournaments() TournamentsRepository
	// Тут можуть бути інші репозиторії, наприклад:
	// Users() UsersRepository
	// Matches() MatchesRepository
}