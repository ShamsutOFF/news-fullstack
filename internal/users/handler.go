package users

import (
	"news-fullstack/internal/users/dto"
	"news-fullstack/pkg/validator"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	router    fiber.Router
	usersRepo *UsersRepository
}

func NewUsersHandler(router fiber.Router, usersRepo *UsersRepository) {
	handler := &UsersHandler{
		router:    router,
		usersRepo: usersRepo,
	}
	api := handler.router.Group("/users")
	api.Post("/register", handler.register)
	api.Post("/login", handler.login)
}

func (h *UsersHandler) register(c *fiber.Ctx) error {
	form := RegisterForm{
		Email:    c.FormValue("email"),
		Name:     c.FormValue("name"),
		Password: c.FormValue("password"),
	}

	// Валидация
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
		&validators.StringLengthInRange{
			Name:    "Password",
			Field:   form.Password,
			Min:     8,
			Max:     100,
			Message: "Пароль должен быть от 8 до 100 символов",
		},
		&validators.StringIsPresent{
			Name:    "Name",
			Field:   form.Name,
			Message: "Имя не задано или не верное",
		},
		&validators.StringLengthInRange{
			Name:    "Name",
			Field:   form.Name,
			Min:     2,
			Max:     100,
			Message: "Имя должно быть от 2 до 100 символов",
		},
	)

	if errors.HasAny() {
		errorMessages := validator.FormatErrors(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Ошибка валидации",
			"details": errorMessages,
		})
	}

	// Проверяем, существует ли пользователь
	exists, err := h.usersRepo.UserExists(form.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Пользователь с таким email уже существует",
		})
	}

	// Хэшируем пароль
	passwordHash, err := hashPassword(form.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при обработке пароля",
		})
	}

	// Создаем пользователя в БД
	err = h.usersRepo.CreateUser(form.Email, form.Name, passwordHash)
	if err != nil {
		if err == ErrUserExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Пользователь с таким email уже существует",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось создать пользователя",
		})
	}

	// По заданию: "При регистрации создавать пользователя и его пока его email просто возвращать в ответе"
	response := dto.RegisterResponse{
		Email: form.Email,
		Name:  form.Name,
	}

	// Если нужно также вернуть ID пользователя (опционально)
	// Можно получить пользователя из БД, чтобы получить его ID
	user, err := h.usersRepo.GetUserByEmail(form.Email)
	if err == nil && user != nil {
		response.ID = user.ID
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *UsersHandler) login(c *fiber.Ctx) error {
	form := RegisterForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	// Валидация
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
	)

	if errors.HasAny() {
		errorMessages := validator.FormatErrors(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Ошибка валидации",
			"details": errorMessages,
		})
	}

	// Получаем пользователя из БД
	user, err := h.usersRepo.GetUserByEmail(form.Email)
	if err != nil {
		if err == ErrUserNotFound {
			// Возвращаем общую ошибку для безопасности
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Неверный email или пароль",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера",
		})
	}

	// Проверяем пароль
	if !checkPassword(form.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неверный email или пароль",
		})
	}

	// Вход успешен
	return c.JSON(fiber.Map{
		"message": "Вход выполнен успешно",
		"user": fiber.Map{
			"email": user.Email,
			"name":  user.Name,
			"id":    user.ID,
		},
	})
}

// Вспомогательные функции для работы с паролями

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
