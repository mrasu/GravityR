package suggest

import (
	"github.com/mrasu/GravityR/database"
	"github.com/rs/zerolog/log"
)

func parseIndexTargets(indexTargetTexts []string) ([]*database.IndexTarget, error) {
	var its []*database.IndexTarget
	for _, text := range indexTargetTexts {
		it, err := database.NewIndexTargetFromText(text)
		if err != nil {
			return nil, err
		}
		its = append(its, it)
	}

	return its, nil
}

func logNewIndexTargets(newIdxTargets []*database.IndexTarget) {
	if len(newIdxTargets) > 0 {
		log.Debug().Msg("Found possibly efficient index combinations:")
		for i, it := range newIdxTargets {
			log.Printf("\t%d.%s", i, it.CombinationString())
		}
	} else {
		log.Debug().Msg("No possibly efficient index found. Perhaps already indexed?")
	}
}
