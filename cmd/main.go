package main

import (
	"log"
	"log/slog"
	"news-fullstack/config"
	"news-fullstack/internal/pages"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2" // ДЛЯ ШАБЛОНОВ!
	slogfiber "github.com/samber/slog-fiber"
)

func main() {
	config.Init()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Инициализация шаблонов - ВАЖНО!
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine, // Подключаем шаблоны к Fiber
	})

	app.Use(slogfiber.New(logger))
	app.Use(recover.New())

	pages.NewPagesHandler(app)

	slog.Info("starting server", "port", 3000)
	log.Fatal(app.Listen(":3000"))
}
