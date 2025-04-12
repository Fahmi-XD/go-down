package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"down/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type imageParams struct {
	Image string `query:"image"`
}

const (
	GROQ string = "https://api.groq.com/openai/v1/chat/completions"
)

var (
	DATA_IMAGE *helper.Visitors = helper.NewVisitors()
)

func init() {
	NewRegister.RegisterProvider(RegisterComponent{
		Title:       "LLM OCR",
		Endpoint:    "/ocr",
		Method:      "GET",
		Description: "Mengekstrak text pada sebuah gambar.",
		Params: map[string]interface{}{
			"image": "https://ytlarge.com/youtube/monetization-checker/captcha",
		},
		Type: "",
		Body: map[string]interface{}{},

		Code: func(c *fiber.Ctx) error {
			params := new(imageParams)

			if err := c.QueryParser(params); err != nil {
				return c.Status(200).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan query image!",
				})
			}

			if params.Image == "" {
				return c.Status(200).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan query image!",
				})
			}

			imageData := saveImage(params.Image)
			js := ocr(c, imageData)

			return c.Status(200).JSON(js)
		},
	})

	NewRegister.RegisterProvider(RegisterComponent{
		Title:       "LLM OCR ( Plugins )",
		Endpoint:    "/ocr-phpsid",
		Method:      "GET",
		Description: "Mengekstrak text pada sebuah gambar.",
		Params: map[string]interface{}{
			"image": "https://ytlarge.com/youtube/monetization-checker/captcha",
		},
		Type:   "",
		Body:   map[string]interface{}{},
		Hidden: true,

		Code: func(c *fiber.Ctx) error {
			params := new(imageParams)

			if err := c.QueryParser(params); err != nil {
				return c.Status(200).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan query image!",
				})
			}

			if params.Image == "" {
				return c.Status(200).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan query image!",
				})
			}

			js := saveImage(params.Image)

			return c.Status(200).JSON(js)
		},
	})
}

func saveImage(image string) map[string]interface{} {
	uid := uuid.New().String()
	req, _ := helper.Request(image, "GET", nil, nil)
	defer req.Body.Close()
	ib, _ := io.ReadAll(req.Body)

	DATA_IMAGE.Write(fmt.Sprintf("image%s.jpg", uid), map[string]any{
		"bytes":       ib,
		"contentType": "image/jpeg",
	})

	go func() { // Menghapus setelah 15 Detik
		time.Sleep(15 * time.Second)

		DATA_IMAGE.Delete(fmt.Sprintf("image%s.jpg", uid))
	}()

	return map[string]interface{}{
		"name": fmt.Sprintf("%s.jpg", uid),
		"uuid": uid,
		"phpsid": req.Header.Values("Set-Cookie"),
	}
}

func ocr(c *fiber.Ctx, imageData map[string]interface{}) map[string]interface{} {
	payload := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "Sebutkan bilangan didalam gambar ini tanpa tambahan kata apapun! Hanya boleh menjawab menggunakan angka contoh: 123, 3466, 457, 2356, 23, 435!!",
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": fmt.Sprintf("%s://%s/api/v1/get-image-nocache?name=%s", c.Protocol(), c.Hostname(), imageData["name"].(string)),
						},
					},
				},
			},
		},
		"model":                 "llama-3.2-90b-vision-preview",
		"temperature":           0.5,
		"max_completion_tokens": 500,
		"top_p":                 0.7,
		"stream":                false,
		"stop":                  nil,
	}

	byt, _ := json.Marshal(payload)

	head := http.Header{}
	head.Set("Content-Type", "application/json")
	head.Set("Authorization", "Bearer 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")

	req, err := helper.Request(GROQ, "POST", bytes.NewReader(byt), head)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()

	resv, _ := io.ReadAll(req.Body)
	var jsn map[string]interface{}
	_ = json.Unmarshal(resv, &jsn)

	jsn["name"] = imageData["name"]
	jsn["uuid"] = imageData["uuid"]
	jsn["phpsid"] = imageData["phpsid"]

	return jsn
}
