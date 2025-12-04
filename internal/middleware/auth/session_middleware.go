package auth

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// SessionMiddleware проверяет авторизацию через сессию
func SessionMiddleware(store *session.Store, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Сохраняем store в контексте для использования в других хендлерах
		c.Locals("session_store", store)

		// Получаем сессию с обработкой ошибок
		sess, err := store.Get(c)
		if err != nil {
			// Логируем ошибку с контекстом запроса
			logger.Error("ошибка получения сессии",
				"error", err,
				"path", c.Path(),
				"method", c.Method(),
				"ip", c.IP(),
			)

			// Устанавливаем флаг ошибки в контексте
			c.Locals("session_error", true)
			c.Locals("session_error_message", err.Error())

			return c.Next() // Продолжаем без сессии
		}

		// Получаем email из сессии
		if email := sess.Get("email"); email != nil {
			// Проверяем тип email перед сохранением
			if emailStr, ok := email.(string); ok {
				c.Locals("user_email", emailStr)
			} else {
				logger.Warn("email в сессии не является строкой",
					"type", email,
					"path", c.Path(),
				)
			}
		}

		// Получаем user_id из сессии (если нужно)
		if userID := sess.Get("user_id"); userID != nil {
			c.Locals("user_id", userID)
		}

		// Получаем name из сессии (если нужно)
		if name := sess.Get("name"); name != nil {
			if nameStr, ok := name.(string); ok {
				c.Locals("user_name", nameStr)
			}
		}

		// Сохраняем сессию в контексте для хендлеров
		c.Locals("session", sess)

		return c.Next()
	}
}
