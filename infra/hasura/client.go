package hasura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	graphqlPath = "/v1/graphql"
	explainPath = "/v1/graphql/explain"
	queryPath   = "/v2/query"
)

type Client struct {
	baseURL *url.URL
	headers map[string]string
}

func NewClient(cfg *Config) *Client {
	return &Client{
		baseURL: cfg.Url,
		headers: map[string]string{
			"x-hasura-admin-secret": cfg.AdminSecret,
		},
	}
}

func (c *Client) QueryWithoutResult(query string, variables map[string]interface{}) error {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	_, err := c.post(graphqlPath, body)
	if err != nil {
		return errors.Wrap(err, "query for GraphQL failed")
	}
	return nil
}

func (c *Client) Explain(body *ExplainRequestBody) ([]*ExplainResponse, error) {
	res, err := c.post(explainPath, body)
	if err != nil {
		return nil, err
	}
	var eres []*ExplainResponse
	err = json.Unmarshal(res, &eres)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response from Hasura")
	}
	return eres, nil
}

func (c *Client) RunRawSQL(sql string) (*RunSQLResponse, error) {
	return c.RunSQL(&RunSQLArgs{
		SQL: sql,
	})
}

func (c *Client) RunSQL(args *RunSQLArgs) (*RunSQLResponse, error) {
	body := map[string]interface{}{
		"type": "run_sql",
		"args": args,
	}
	res, err := c.post(queryPath, body)
	if err != nil {
		return nil, err
	}
	rres := &RunSQLResponse{}
	err = json.Unmarshal(res, rres)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from Hasura")
	}
	return rres, nil
}

func (c *Client) GetTableColumns(schema string, tableNames []string) ([]*ColumnInfo, error) {
	q, err := buildColumnFetchQuery(schema, tableNames)
	if err != nil {
		return nil, err
	}

	args := &RunSQLArgs{
		SQL:      q,
		ReadOnly: lib.Ptr(true),
	}
	res, err := c.RunSQL(args)
	if err != nil {
		return nil, err
	}
	cols := parseColumnFetchResult(res)
	return cols, nil
}

func (c *Client) post(pathname string, body interface{}) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}
	req, err := http.NewRequest("POST", c.resolveURL(pathname), bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create POST request to %s", pathname))
	}
	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute http request to hasura")
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		text, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Request to Hasura failed (Status=%d)", res.StatusCode))
		}

		return nil, errors.Errorf("Request to Hasura failed (Status=%d): %s", res.StatusCode, string(text))
	}

	return ioutil.ReadAll(res.Body)
}

func (c *Client) resolveURL(pathname string) string {
	u := *c.baseURL
	u.Path = path.Join(u.Path, pathname)
	return u.String()
}
