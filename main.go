package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func indexHandler(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func main() {
	fmt.Println("PostgreSQL Multi-Partitions")
	app := fiber.New()
	app.Get("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
