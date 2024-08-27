package config

import (
	"fmt"
	"gofi/pkg/constant"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/masb0ymas/go-utils/pkg"
)

func Cors() cors.Config {
	allowedOrigin := strings.Join(constant.AllowedOrigin(), ", ")

	logMessage := pkg.PrintLog("Cors", "Allowed Origins ( "+allowedOrigin+" )")
	fmt.Println(logMessage)

	result := cors.Config{
		AllowOrigins: allowedOrigin,
		// AllowMethods:  "GET, POST, HEAD, PUT, DELETE, PATCH",
		// AllowHeaders:  "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token",
		// ExposeHeaders: "Content-Length",
		// MaxAge:        86400,
	}

	return result
}
