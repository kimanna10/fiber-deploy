package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Драйвер Postgres
)

func main() {
	// 1. Загружаем переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}

	// 2. Строка подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// 3. Подключение к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	// 4. Пример маршрута с обращением к БД
	app.Get("/db-test", func(c *fiber.Ctx) error {
		var now string
		// Простой запрос, чтобы проверить связь
		err := db.QueryRow("SELECT NOW()").Scan(&now)
		if err != nil {
			return c.Status(500).SendString("Ошибка БД: " + err.Error())
		}
		return c.SendString("Время сервера БД: " + now)
	})

	log.Fatal(app.Listen(":3000"))
}
