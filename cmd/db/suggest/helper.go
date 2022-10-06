package suggest

import (
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/rs/zerolog/log"
)

func parseIndexTargets(indexTargetTexts []string) ([]*dmodel.IndexTarget, error) {
	var its []*dmodel.IndexTarget
	for _, text := range indexTargetTexts {
		it, err := dmodel.NewIndexTargetFromText(text)
		if err != nil {
			return nil, err
		}
		its = append(its, it)
	}

	return its, nil
}

func logNewIndexTargets(newIdxTargets []*dmodel.IndexTarget) {
	if len(newIdxTargets) > 0 {
		log.Debug().Msg("Found possibly efficient index combinations:")
		for i, it := range newIdxTargets {
			log.Printf("\t%d.%s", i, it.CombinationString())
		}
	} else {
		log.Debug().Msg("No possibly efficient index found. Perhaps already indexed?")
	}
}
