package session

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSessionStore(dbpool *pgxpool.Pool) *session.Store {
	// Создаем хранилище в Postgres
	storage := postgres.New(postgres.Config{
		ConnectionURI: "", // оставляем пустым, используем dbpool
		DB:            dbpool,
		Table:         "sessions",    // таблица для сессий
		Reset:         false,         // не очищать таблицу при запуске
		GCInterval:    1 * time.Hour, // очистка старых сессий каждый час
	})

	// Создаем store для сессий
	store := session.New(session.Config{
		Storage:        storage,
		Expiration:     24 * time.Hour,      // сессия живет 24 часа
		KeyLookup:      "cookie:session_id", // кука будет называться session_id
		CookieHTTPOnly: true,                // защита от XSS
		CookieSecure:   false,               // true для продакшена (HTTPS)
		CookieSameSite: "Lax",
	})

	return store
}
