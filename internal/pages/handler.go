package pages

import (
	"github.com/gofiber/fiber/v2"
	"news-fullstack/pkg/tadapter"
	"news-fullstack/views"
)

type PagesHandler struct {
	router fiber.Router
}

func NewPagesHandler(router fiber.Router) {
	handler := &PagesHandler{
		router: router,
	}
	api := handler.router.Group("/api")
	api.Get("/", handler.home)
	api.Get("/error", handler.error)
}

func (h *PagesHandler) home(c *fiber.Ctx) error {
	//categories := []string{
	//	"🍕 Еда",
	//	"⚽ Спорт",
	//	"🚗 Машины",
	//	"🐶 Животные",
	//	"💻 Технологии",
	//	"🎬 Фильмы",
	//	"🎵 Музыка",
	//	"🌆 Путешествия",
	//	"💼 Бизнес",
	//	"🏥 Здоровье",
	//}
	component := views.Hello("Adel")
	return tadapter.Render(c, component)
}

func (h *PagesHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
