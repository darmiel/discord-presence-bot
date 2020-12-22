package main

import (
	"github.com/darmiel/discord-presence-bot/internal"
)

func init() {
	internal.InitFlags()
	internal.InitDotEnv()

	internal.InitDatabase()
}

func main() {
	internal.CreateBot()
}
