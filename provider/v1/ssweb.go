package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
)

type sswebRequest struct {
	Device string `query:"device"`
	URL    string `query:"url"`
}

var (
	deviceMap map[string][]int = map[string][]int{
		"mobile-small":  {320, 568},
		"mobile-medium": {375, 667},
		"mobile-large":  {425, 736},
		"tablet":        {768, 1024},
		"laptop":        {1024, 1366},
		"desktop":       {1440, 900},
		"large-screen":  {1920, 1080},
	}
)

func init() {
	NewRegister.RegisterProvider(RegisterComponent{
		Title:       "Screenshot website ( List Device )",
		Endpoint:    "/ssweb-device",
		Method:      "GET",
		Description: "Melihat daftar device yang tersedia",
		Params:      map[string]interface{}{},
		Type:        "",
		Body:        map[string]interface{}{},

		Code: func(c *fiber.Ctx) error {
			return c.Status(200).JSON(fiber.Map{
				"device": []string{
					"mobile-small",
					"mobile-medium",
					"mobile-large",
					"tablet",
					"laptop",
					"desktop",
					"large-screen",
				},
			})
		},
	})

	NewRegister.RegisterProvider(RegisterComponent{
		Title:       "Screenshot website",
		Endpoint:    "/ssweb",
		Method:      "GET",
		Description: "Screenshot halaman website",
		Params: map[string]interface{}{
			"device": "mobile",
		},
		Type: "",
		Body: map[string]interface{}{},

		Code: func(c *fiber.Ctx) error {
			query := new(sswebRequest)
			if err := c.QueryParser(query); err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan query 'device' dan 'url'",
					"device": []string{
						"mobile-small",
						"mobile-medium",
						"mobile-large",
						"tablet",
						"laptop",
						"desktop",
						"large-screen",
					},
				})
			}
			if query.Device == "" || !isAvailable(query.Device) || query.URL == "" {
				return c.Status(400).JSON(fiber.Map{
					"error":   true,
					"message": "Masukan query 'device' dan 'url'",
					"device": []string{
						"mobile-small",
						"mobile-medium",
						"mobile-large",
						"tablet",
						"laptop",
						"desktop",
						"large-screen",
					},
				})
			}

			result := sweb(query.Device, query.URL)

			c.Response().Header.Set("Content-Type", "image/png")
			return c.Status(200).Send(result)
		},
	})
}

func isAvailable(device string) bool {
	for k := range deviceMap {
		if k == device {
			return true
		}
	}

	return false
}

func isMobile(device string) bool {
	if strings.HasPrefix(device, "mobile") {
		return true
	}

	return false
}

func sweb(device, link string) []byte {
	fmt.Println("[ SSWEB ] Mengambil screenshot...")

	ctx, cancel, err := cu.New(cu.NewConfig(
		cu.WithHeadless(),
		cu.WithTimeout(30*time.Second),
	))
	defer cancel()
	if err != nil {
		fmt.Println(err)
	}

	var result []byte
	if err := chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(int64(deviceMap[device][0]), int64(deviceMap[device][1]), 1.0, isMobile(device)),
		chromedp.Navigate(link),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			capture, err := page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  float64(int64(deviceMap[device][0])),
					Height: float64(int64(deviceMap[device][1])),
					Scale:  1,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			result = capture
			return nil
		}),
	); err != nil {
		fmt.Println(err)
	}

	return result
}
