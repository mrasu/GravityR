package util

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

func LogError(err error) {
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		log.Error().Msgf("%+v", err)
	} else {
		log.Error().Msg(err.Error())
	}
}

func LogResultOutputPath(outputPath string) {
	wd, err := os.Getwd()
	if err == nil {
		log.Info().Msg("Result html is at: " + path.Join(wd, outputPath))
	}
}
