package pages

import (
	"news-fullstack/pkg/tadapter"
	"news-fullstack/views"
	"news-fullstack/views/types"

	"github.com/gofiber/fiber/v2"
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
		{Name: "Музыка", ImageURL: "/public/images/music_img.jpg"},
		{Name: "Технологии", ImageURL: "/public/images/tech_img.jpg"},
		{Name: "Прочее", ImageURL: "/public/images/other_img.jpg"},
	}
	banners := []types.BannerCard{
		{
			ImageURL:    "/public/images/car_big_img.jpg",
			Title:       "Как безопасно водить",
			Description: "Длинный текст про то, как можно безопасно водить автомобиль.",
		},
		{
			ImageURL:    "/public/images/music_big_img.jpg",
			Title:       "Создавай музыку!",
			Description: "Сегодня мы рассмотрим технику быстрого создания музыки за счёт использования...",
		},
	}

	component := views.HomePage(categories, banners)
	return tadapter.Render(c, component)
}

func (h *PagesHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
