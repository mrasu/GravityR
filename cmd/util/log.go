package util

import (
	"github.com/mrasu/GravityR/database/dmodel"
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

func LogNewIndexTargets(newIdxTargets []*dmodel.IndexTarget) {
	if len(newIdxTargets) > 0 {
		log.Debug().Msg("Found possibly efficient index combinations:")
		for i, it := range newIdxTargets {
			log.Printf("\t%d.%s", i, it.CombinationString())
		}
	} else {
		log.Debug().Msg("No possibly efficient index found. Perhaps already indexed?")
	}
}
