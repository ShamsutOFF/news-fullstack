package api

import (
	"news-fullstack/pkg/validator"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
)

type ApiHandler struct {
	router fiber.Router
}

func NewApiHandler(router fiber.Router) {
	handler := &ApiHandler{
		router: router,
	}
	api := handler.router.Group("/api")
	api.Post("/register", handler.register)
	api.Get("/error", handler.error)
}

func (h *ApiHandler) register(c *fiber.Ctx) error {
	form := RegisterForm{
		Email:    c.FormValue("email"),
		Name:     c.FormValue("name"),
		Password: c.FormValue("password"),
	}

	errors := validate.Validate(
		&validators.EmailIsPresent{
			Name:    "Email",
			Field:   form.Email,
			Message: "Email не задан или не верный",
		},
		&validators.StringIsPresent{
			Name:    "Password",
			Field:   form.Password,
			Message: "Пароль не задан или не верный",
		},
		&validators.StringIsPresent{
			Name:    "Name",
			Field:   form.Name,
			Message: "Имя не задано или не верное",
		},
	)
	result := "Вы зарегистрировадись!"
	if errors.HasAny() {
		result = validator.FormatErrors(errors)
	}
	return c.SendString(result)
}

func (h *ApiHandler) error(c *fiber.Ctx) error {
	return c.SendString("Hello, World from Error 👋!")
}
