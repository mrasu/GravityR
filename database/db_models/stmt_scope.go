package db_models

import "github.com/mrasu/GravityR/lib"

const RootScopeName = "<root>"

type StmtScope struct {
	Name        string
	Fields      []*Field
	FieldScopes []*StmtScope
	Tables      []*Table
	Scopes      []*StmtScope
	CTEs        map[string]*StmtScope
}

func (ss *StmtScope) ListAsTableMap() map[string]*lib.Set[string] {
	res := map[string]*lib.Set[string]{}
	for _, t := range ss.Tables {
		res[t.AsName] = lib.NewSetS([]string{t.Name})
	}

	for _, s := range ss.Scopes {
		asTableMap := s.ListAsTableMap()
		ts := lib.NewSet[string]()
		for name, tNames := range asTableMap {
			//TODO: consider scope to handle name duplication
			if s, ok := res[name]; ok {
				s.Merge(tNames)
			} else {
				res[name] = lib.NewSetS(tNames.Values())
			}

			ts.Merge(tNames)
		}
		res[s.Name] = ts
	}

	delete(res, "")
	return res
}
