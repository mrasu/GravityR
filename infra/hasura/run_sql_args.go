package hasura

type RunSQLArgs struct {
	SQL                      string  `json:"sql"`
	Source                   *string `json:"source"`
	Cascade                  *bool   `json:"cascade"`
	CheckMetadataConsistency *bool   `json:"check_metadata_consistency"`
	ReadOnly                 *bool   `json:"read_only"`
}
