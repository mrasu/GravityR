package service

import (
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/lib"
)

// BuildExplainedIndexTargets extracts IndexTargetTable referred in EXPLAIN ANALYZE(aka SingleTableExplainResult) and orders them by its EstimatedTotalTime
func BuildExplainedIndexTargets(idxCandidates []*database.IndexTargetTable, scopes []*database.StmtScope, explainResults []*database.SingleTableExplainResult) ([]*database.IndexTargetTable, []error) {
	//TODO: consider scope to handle name duplication
	asTableMap := map[string]*lib.Set[string]{}
	for _, s := range scopes {
		for tName, tables := range s.ListAsTableMap() {
			if _, ok := asTableMap[tName]; ok {
				asTableMap[tName].Merge(tables)
			} else {
				asTableMap[tName] = tables
			}
		}
	}
	lib.Sort(explainResults, func(t *database.SingleTableExplainResult) float64 {
		return t.EstimatedTotalTime * -1
	})

	calledTable := lib.NewSet[string]()
	var indexes []*database.IndexTargetTable
	for _, eRes := range explainResults {
		tNames := []string{eRes.TableName}
		if name, ok := asTableMap[eRes.TableName]; ok {
			tNames = name.Values()
		}

		for _, tName := range tNames {
			if calledTable.Contains(tName) {
				continue
			}
			for _, c := range idxCandidates {
				if c.TableName == tName {
					indexes = append(indexes, c)
				}
			}
			calledTable.Add(tName)
		}
	}

	return indexes, nil
}
