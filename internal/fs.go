package internal

import (
	"embed"
)

//go:embed migration/*.sql
var EmbedMigration embed.FS
