module bunny-bot

go 1.24.3

tool (
	bunny-bot/cmd/bunny-bot
	bunny-bot/internal/moon
)

require (
	github.com/NicoNex/echotron/v3 v3.46.0
	github.com/akyoto/cache v1.0.6
	github.com/davecgh/go-spew v1.1.1
	github.com/hjson/hjson-go v3.3.0+incompatible
)

require golang.org/x/time v0.12.0 // indirect
