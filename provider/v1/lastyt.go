package v1

import (
	"down/helper"

	"github.com/gofiber/fiber/v2"
)

type lastytQuery struct {
	ID string `query:"id"`
}

var LS *helper.Visitors = helper.NewVisitors()

func init() {
	NewRegister.RegisterProvider(RegisterComponent{
		Title:       "Last youtube provider ( Set )",
		Endpoint:    "/set-lastyt",
		Method:      "GET",
		Description: "Menset informasi tentang provider terakhir yang aktif",
		Params: map[string]interface{}{
			"id": "1",
		},
		Type:   "",
		Body:   map[string]interface{}{},
		Hidden: true,

		Code: func(c *fiber.Ctx) error {
			query := new(lastytQuery)

			if err := c.QueryParser(query); err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan id yang ingin di set!",
				})
			}

			if query.ID == "" {
				return c.Status(400).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan id yang ingin di set!",
				})
			}

			LS.Write("lastyt", query.ID)

			return c.Status(200).JSON(fiber.Map{
				"status": "ok",
				"current-id": query.ID,
			})
		},
	})

	NewRegister.RegisterProvider(RegisterComponent{
		Title:       "Last youtube provider",
		Endpoint:    "/get-lastyt",
		Method:      "GET",
		Description: "Mendapatkan informasi tentang provider terakhir yang aktif",
		Params:      map[string]interface{}{},
		Type:        "",
		Body:        map[string]interface{}{},
		Hidden:      true,

		Code: func(c *fiber.Ctx) error {
			lastId := LS.Read("lastyt")

			return c.Status(200).JSON(fiber.Map{
				"status": "ok",
				"current-id": lastId,
			})
		},
	})
}
