package main

import (
	"golang-fiber-chatapp/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/", "./public")

	routes.SetupChatRoutes(app)

	log.Println("Chat server is running on ws://localhost:3000")
	app.Listen(":3000")
}
