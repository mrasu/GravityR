package suggest

import (
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/lib"
)

func parseIndexTargets(indexTargetTexts []string) ([]*common_model.IndexTarget, error) {
	var its []*common_model.IndexTarget
	for _, text := range indexTargetTexts {
		it, err := common_model.NewIndexTarget(text)
		if err != nil {
			return nil, err
		}
		its = append(its, it)
	}

	return its, nil
}

func toUniqueIndexTargets(its []*common_model.IndexTargetTable) []*common_model.IndexTarget {
	var idxTargets []*common_model.IndexTarget
	for _, it := range its {
		idxTargets = append(idxTargets, it.ToIndexTarget())
	}

	return lib.BruteUniq(idxTargets)
}
