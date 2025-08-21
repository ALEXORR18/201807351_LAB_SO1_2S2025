package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Mensaje string `json:"mensaje"`
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func selfMessage() Message {
	api := env("API_NAME", "API2")
	vm := env("VM_NAME", "VM1")
	student := env("STUDENT_NAME", "Brian Alexander García Orr")
	carnet := env("CARNET", "201807351")
	return Message{
		Mensaje: "Hola, responde la API: " + api + " en la " + vm +
			", desarrollada por el estudiante " + student + " con carnet: " + carnet,
	}
}

func main() {
	app := fiber.New()
	client := &http.Client{Timeout: 4 * time.Second}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(selfMessage())
	})

	// /api2/201807351/llamar-api1
	app.Get("/api2/201807351/llamar-api1", func(c *fiber.Ctx) error {
		url := env("API1_URL", "http://VM1:3000")
		resp, err := client.Get(url)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error llamando a API1", "detalle": err.Error()})
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return c.Send(body)
	})

	// /api2/201807351/llamar-api3
	app.Get("/api2/201807351/llamar-api3", func(c *fiber.Ctx) error {
		url := env("API3_URL", "http://VM2:3002")
		resp, err := client.Get(url)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error llamando a API3", "detalle": err.Error()})
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return c.Send(body)
	})

	app.Listen("0.0.0.0:3001")
}
