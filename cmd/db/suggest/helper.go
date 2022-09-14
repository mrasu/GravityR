package suggest

import (
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/lib"
)

func parseIndexTargets(indexTargetTexts []string) ([]*db_models.IndexTarget, error) {
	var its []*db_models.IndexTarget
	for _, text := range indexTargetTexts {
		it, err := db_models.NewIndexTarget(text)
		if err != nil {
			return nil, err
		}
		its = append(its, it)
	}

	return its, nil
}

func toUniqueIndexTargets(its []*db_models.IndexTargetTable) []*db_models.IndexTarget {
	var idxTargets []*db_models.IndexTarget
	for _, it := range its {
		idxTargets = append(idxTargets, it.ToIndexTarget())
	}

	return lib.BruteUniq(idxTargets)
}
