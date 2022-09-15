package hasura

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/lib"
	"time"
)

type IndexExaminer struct {
	cli       *hasura.Client
	query     string
	variables map[string]interface{}
}

func NewIndexExaminer(cli *hasura.Client, query string, v map[string]interface{}) *IndexExaminer {
	return &IndexExaminer{
		cli:       cli,
		query:     query,
		variables: v,
	}
}

func (ie *IndexExaminer) Execute() (int64, error) {
	start := time.Now()
	err := ie.cli.QueryWithoutResult(ie.query, ie.variables)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed.Milliseconds(), nil
}

func (ie *IndexExaminer) CreateIndex(name string, it *common_model.IndexTarget) error {
	sql := fmt.Sprintf(`CREATE INDEX "%s" ON "%s" (%s)`,
		name, it.TableName,
		lib.Join(it.Columns, ",", func(i *common_model.IndexColumn) string { return `"` + i.SafeName() + `"` }),
	)
	_, err := ie.cli.RunRawSQL(sql)
	return err
}

func (ie *IndexExaminer) DropIndex(name string, it *common_model.IndexTarget) error {
	sql := fmt.Sprintf(`DROP INDEX "%s"`, name)
	_, err := ie.cli.RunRawSQL(sql)
	return err
}
