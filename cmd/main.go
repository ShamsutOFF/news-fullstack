package main

import (
	"github.com/gofiber/fiber/v3"
	recover2 "github.com/gofiber/fiber/v3/middleware/recover"
	"log"
	"news-fullstack/config"
	"news-fullstack/internal/pages"
)

func main() {
	config.Init()
	dbConf := config.NewDatabaseConfig()
	log.Println(dbConf)

	app := fiber.New()
	app.Use(recover2.New())

	pages.NewPagesHandler(app)

	log.Fatal(app.Listen(":3000"))
}
