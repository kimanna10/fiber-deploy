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
	// 1. Загружаем .env (игнорируем ошибку, если файла нет — берем из системы)
	godotenv.Load()

	// 2. Формируем строку подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// 3. Подключение к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка конфигурации БД: %v", err)
	}
	defer db.Close()

	// 4. Проверка связи с БД (Ping)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL!")

	// 5. Инициализация Fiber
	app := fiber.New()

	// 6. Роуты
	app.Get("/db-test", func(c *fiber.Ctx) error {
		var now string
		err := db.QueryRow("SELECT NOW()").Scan(&now)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Ошибка запроса к БД: " + err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"time":   now,
		})
	})

	// 7. Запуск сервера
	log.Fatal(app.Listen(":3000"))
}
