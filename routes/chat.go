// routes/chat.go
package routes

import (
	"golang-fiber-chatapp/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupChatRoutes(app *fiber.App) {

	// Middleware to handle WebSocket upgrades for the "/ws" route
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint
	app.Get("/ws", websocket.New(handlers.ChatHandler))
}
