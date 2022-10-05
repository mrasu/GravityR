package mysql

const indexFetchQuery = `
SELECT
	TABLE_NAME,
	INDEX_NAME,
	COLUMN_NAME
FROM
	information_schema.STATISTICS
WHERE
	TABLE_SCHEMA = ? AND
	TABLE_NAME IN (?)
ORDER BY
	TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
`

type IndexInfo struct {
	TableName string
	Columns   []string
}

type flatIndexInfo struct {
	TableName  string `db:"TABLE_NAME"`
	IndexName  string `db:"INDEX_NAME"`
	ColumnName string `db:"COLUMN_NAME"`
}

func toIndexInfo(orderedFlatInfo []*flatIndexInfo) []*IndexInfo {
	var infos []*IndexInfo

	lastTName := ""
	lastIName := ""
	for _, flatInfo := range orderedFlatInfo {
		if lastTName == flatInfo.TableName && lastIName == flatInfo.IndexName {
			idx := infos[len(infos)-1]
			idx.Columns = append(idx.Columns, flatInfo.ColumnName)
		} else {
			infos = append(infos, &IndexInfo{
				TableName: flatInfo.TableName,
				Columns:   []string{flatInfo.ColumnName},
			})
		}

		lastTName = flatInfo.TableName
		lastIName = flatInfo.IndexName
	}

	return infos
}
