package util

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogError(err error) {
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		log.Error().Msgf("%+v", err)
	} else {
		log.Error().Msg(err.Error())
	}
}
