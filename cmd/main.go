package main

import (
	"log"
	"log/slog"
	"news-fullstack/config"
	"news-fullstack/internal/pages"
	"news-fullstack/internal/users"
	"news-fullstack/pkg/database"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
)

func main() {
	config.Init()
	dbConfig := config.NewDatabaseConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := fiber.New()

	app.Use(slogfiber.New(logger))
	app.Use(recover.New())
	app.Static("/public", "./public")

	dbpool := database.CreateDbPool(dbConfig, logger)
	defer dbpool.Close()

	usersRepo := users.NewUsersRepository(dbpool, logger)
	pages.NewPagesHandler(app, usersRepo)
	users.NewUsersHandler(app, usersRepo)

	slog.Info("starting server", "port", 3000)
	log.Fatal(app.Listen(":3000"))
}
