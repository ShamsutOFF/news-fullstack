package pages

import (
	"github.com/gofiber/fiber/v2"
	"news-fullstack/pkg/tadapter"
	"news-fullstack/views"
	"news-fullstack/views/types"
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
	categories := []types.Category{
		{Name: "Еда", ImageURL: "/public/images/food_img.jpg"},
		{Name: "Животные", ImageURL: "/public/images/animal_img.jpg"},
		{Name: "Машины", ImageURL: "/public/images/car_img.jpg"},
		{Name: "Спорт", ImageURL: "/public/images/sport_img.jpg"},
		{Name: "Музыка", ImageURL: "/public/images/sport_img.jpg"},
		{Name: "Технологии", ImageURL: "/public/images/tech_img.jpg"},
		{Name: "Прочее", ImageURL: "/public/images/other_img.jpg"},
	}

	component := views.HomePage(categories)
	return tadapter.Render(c, component)
}

func (h *PagesHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
