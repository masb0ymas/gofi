// Package migrations embeds the golang-migrate SQL files so both cmd/api
// (migrate on boot) and cmd/migrate run them without a filesystem path.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
