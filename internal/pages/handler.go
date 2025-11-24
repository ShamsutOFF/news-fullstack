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

	articles := []types.ArticleCard{
		{
			ImageURL:     "/public/images/art1_img.jpg",
			Title:        "Открытие сезона байдарок",
			Description:  "Сегодня был открыт сезон путешествия на байдарках, где вы можете поучаствовать в увлекательных маршрутах по живописным озерам и рекам нашего региона.",
			AuthorName:   "Михаил Аршинов",
			AuthorAvatar: "/public/images/author_avatar1.jpg",
			PublishDate:  "Август 18, 2025",
		},
		{
			ImageURL:     "/public/images/art2_img.jpg",
			Title:        "Выбери правильный ноутбук для задач",
			Description:  "От верного выбора ноутбука зависит не только удобство, но и эффективность работы...",
			AuthorName:   "Анна Петрова",
			AuthorAvatar: "/public/images/author_avatar2.jpg",
			PublishDate:  "Август 15, 2025",
		},
		{
			ImageURL:     "/public/images/art3_img.jpg",
			Title:        "Создание автомобилей с автопилотом",
			Description:  "Электические автомобили без водителя скоро станут реальностью, где нам не придётся ...",
			AuthorName:   "Дмитрий Соколов",
			AuthorAvatar: "/public/images/author_avatar3.jpg",
			PublishDate:  "Август 12, 2025",
		},
		{
			ImageURL:     "/public/images/art4_img.jpg",
			Title:        "Как быстро приготовить вкусный обед",
			Description:  "Сегодня поговорим о том, как можно быстро и эффективно приготовить обед для ...",
			AuthorName:   "Дмитрий Соколов",
			AuthorAvatar: "/public/images/author_avatar4.jpg",
			PublishDate:  "Август 12, 2025",
		},
	}

	component := views.HomePage(categories, banners, articles)
	return tadapter.Render(c, component)
}

func (h *PagesHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
