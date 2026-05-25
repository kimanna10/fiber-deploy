package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Println("env не найден, использую переменные окружения")
	}

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("ошибка парсинга конфигурации БД: %v", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	log.Println("Подключение установлено!")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		var now string
		err := db.QueryRow(c.UserContext(), "SELECT NOW()").Scan(&now)
		if err != nil {
			return c.Status(500).SendString("Ошибка запроса к БД")
		}
		return c.JSON(fiber.Map{"current_time": now})
	})

	log.Fatal(app.Listen(":3000"))
}
