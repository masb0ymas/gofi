package middlewares

import (
	"net/url"
	"slices"

	"github.com/gofiber/fiber/v3"

	"gofi/internal/lib/apperr"
)

// crossOrigin reports whether a browser sent a state-changing request from an
// origin other than this API's own. Browsers set Sec-Fetch-Site themselves and
// scripts cannot change it; a request without it comes from a non-browser
// client, or an old browser, and is judged by its Origin header instead.
func crossOrigin(method, fetchSite, origin, host string, trusted []string) bool {
	switch method {
	case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
		return false
	}
	if slices.Contains(trusted, origin) {
		return false
	}

	switch fetchSite {
	case "same-origin", "none":
		return false
	case "":
		if origin == "" {
			return false
		}
		u, err := url.Parse(origin)
		return err != nil || u.Host != host
	}
	// cross-site, and same-site: a sibling subdomain is not trusted either.
	return true
}

// CrossOrigin refuses state-changing requests that a browser sends from
// another origin, unless that origin is trusted. It is the CSRF defence for
// the cookie session, and needs no token from clients.
func CrossOrigin(trusted []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if crossOrigin(c.Method(), c.Get("Sec-Fetch-Site"), c.Get(fiber.HeaderOrigin), c.Host(), trusted) {
			return apperr.New(apperr.KindForbidden, "cross-origin request refused")
		}
		return c.Next()
	}
}
