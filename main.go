package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env не найден, использую переменные окружения")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("переменная не задана")
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	log.Println("Подключение установлено!")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		var now string
		err := db.QueryRow("SELECT NOW()").Scan(&now)
		if err != nil {
			return c.Status(500).SendString("Ошибка запроса к БД")
		}
		return c.JSON(fiber.Map{"current_time": now})
	})

	log.Fatal(app.Listen(":3000"))
}
