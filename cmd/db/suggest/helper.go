package suggest

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/rs/zerolog/log"
)

func parseIndexTargets(indexTargetTexts []string) ([]*common_model.IndexTarget, error) {
	var its []*common_model.IndexTarget
	for _, text := range indexTargetTexts {
		it, err := common_model.NewIndexTargetFromText(text)
		if err != nil {
			return nil, err
		}
		its = append(its, it)
	}

	return its, nil
}

func logNewIndexTargets(newIdxTargets []*common_model.IndexTarget) {
	if len(newIdxTargets) > 0 {
		log.Debug().Msg("Found possibly efficient index combinations:")
		for i, it := range newIdxTargets {
			log.Printf("\t%d.%s", i, it.CombinationString())
		}
	} else {
		log.Debug().Msg("No possibly efficient index found. Perhaps already indexed?")
	}
}
