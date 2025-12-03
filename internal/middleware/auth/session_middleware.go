package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// SessionMiddleware проверяет авторизацию через сессию
func SessionMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Сохраняем store в контексте для использования в других хендлерах
		c.Locals("session_store", store)

		// Получаем сессию
		sess, err := store.Get(c)
		if err != nil {
			return c.Next() // Продолжаем без сессии
		}

		// Получаем email из сессии
		email := sess.Get("email")
		if email != nil {
			// Сохраняем email в контексте для использования в шаблонах
			c.Locals("user_email", email)
		}

		// Сохраняем сессию в контексте для хендлеров
		c.Locals("session", sess)
		return c.Next()
	}
}
