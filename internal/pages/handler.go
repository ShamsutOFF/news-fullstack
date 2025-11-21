package pages

import "github.com/gofiber/fiber/v2"

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
	return c.SendString("Hello, World from Home 👋!")
}

func (h *PagesHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
