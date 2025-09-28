package config

import (
	"gofi/src/lib/constant"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/masb0ymas/go-utils/pkg"
)

func Cors() cors.Config {
	allowedOrigin := strings.Join(constant.AllowedOrigin(), ", ")

	msg := pkg.Println("Cors", "Allowed Origins ( "+allowedOrigin+" )")
	log.Println(msg)

	result := cors.Config{
		AllowOrigins: allowedOrigin,
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		// AllowHeaders:  "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token",
		// ExposeHeaders: "Content-Length",
		// MaxAge:        86400,
	}

	return result
}
