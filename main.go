package main

import (

	"os"

	"goAPI/v2/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	api.StartApi()
}

