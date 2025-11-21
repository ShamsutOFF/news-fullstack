package pages

import (
	"github.com/gofiber/fiber/v2"
)

type PagesHandler struct {
	router fiber.Router
}

// Данные для шаблона
type TemplateData struct {
	Categories []string
	Message    string
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
	// Создаем список категорий
	categories := []string{
		"🍕 Еда",
		"⚽ Спорт",
		"🚗 Машины",
		"🐶 Животные",
		"💻 Технологии",
		"🎬 Фильмы",
		"🎵 Музыка",
		"🌆 Путешествия",
		"💼 Бизнес",
		"🏥 Здоровье",
	}

	// Подготовка данных для шаблона
	data := TemplateData{
		Categories: categories,
		Message:    "Добро пожаловать на наш новостной портал! Выберите категорию выше.",
	}

	// Рендерим шаблон с данными
	return c.Render("home", data)
}

func (h *PagesHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
