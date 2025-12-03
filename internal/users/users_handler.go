package users

import (
	"fmt"
	"log"
	"news-fullstack/pkg/validator"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	api.Get("/logout", handler.logout)
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
		return c.SendString(fmt.Sprintf("❌ Ошибка валидации: %s", errorMessages))
	}

	// Проверяем, существует ли пользователь
	exists, err := h.usersRepo.UserExists(form.Email)
	if err != nil {
		return c.SendString("❌ Внутренняя ошибка сервера")
	}
	if exists {
		return c.SendString("❌ Пользователь с таким email уже существует")
	}

	// Хэшируем пароль
	passwordHash, err := hashPassword(form.Password)
	if err != nil {
		return c.SendString("❌ Ошибка при обработке пароля")
	}

	// Создаем пользователя в БД
	err = h.usersRepo.CreateUser(form.Email, form.Name, passwordHash)
	if err != nil {
		if err == ErrUserExists {
			return c.SendString("❌ Пользователь с таким email уже существует")
		}
		return c.SendString("❌ Не удалось создать пользователя")
	}

	// По заданию: "При регистрации создавать пользователя и его пока его email просто возвращать в ответе"
	// Форматируем красивый ответ
	response := fmt.Sprintf("✅ Регистрация успешна!\n\nEmail: %s\nИмя: %s", form.Email, form.Name)

	// Если нужно также вернуть ID пользователя (опционально)
	// Можно получить пользователя из БД, чтобы получить его ID
	user, err := h.usersRepo.GetUserByEmail(form.Email)
	if err == nil && user != nil {
		response = fmt.Sprintf("✅ Регистрация успешна!\n\nEmail: %s\nИмя: %s\nID: %d",
			form.Email, form.Name, user.ID)
	}

	return c.SendString(response)
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
		return c.SendString(fmt.Sprintf("❌ Ошибка валидации: %s", errorMessages))
	}

	// Получаем пользователя из БД
	user, err := h.usersRepo.GetUserByEmail(form.Email)
	if err != nil {
		if err == ErrUserNotFound {
			// Возвращаем общую ошибку для безопасности
			return c.SendString("❌ Неверный email или пароль")
		}
		return c.SendString("❌ Внутренняя ошибка сервера")
	}

	// Проверяем пароль
	if !checkPassword(form.Password, user.PasswordHash) {
		return c.SendString("❌ Неверный email или пароль")
	}

	// Создаем сессию
	sessionStore := c.Locals("session_store").(*session.Store)
	sess, err := sessionStore.Get(c)
	if err != nil {
		log.Printf("Ошибка получения сессии: %v", err)
		return c.Status(500).SendString("❌ Ошибка создания сессии")
	}

	// Сохраняем данные пользователя в сессии
	sess.Set("email", user.Email)
	sess.Set("user_id", user.ID)
	sess.Set("name", user.Name)

	// Сохраняем сессию
	if err := sess.Save(); err != nil {
		log.Printf("Ошибка сохранения сессии: %v", err)
		return c.Status(500).SendString("❌ Ошибка сохранения сессии")
	}

	// Редирект на страницу успеха
	return c.Redirect(fmt.Sprintf("/login-success?name=%s", user.Name))
}

// Добавим метод для выхода
func (h *UsersHandler) logout(c *fiber.Ctx) error {
	if storeInterface := c.Locals("session_store"); storeInterface != nil {
		if store, ok := storeInterface.(*session.Store); ok {
			sess, err := store.Get(c)
			if err == nil {
				sess.Destroy()
			}
		}
	}
	return c.Redirect("/")
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
