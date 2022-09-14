package hasura

import (
	"fmt"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"net/url"
)

type Config struct {
	Url         *url.URL
	AdminSecret string
}

func NewConfigFromEnv() (*Config, error) {
	u, err := lib.GetEnv("HASURA_URL")
	if err != nil {
		return nil, err
	}
	uu, err := url.Parse(u)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot parse URL for Hasura: %s", u))
	}

	s, err := lib.GetEnv("HASURA_ADMIN_SECRET")
	if err != nil {
		return nil, err
	}

	return &Config{
		Url:         uu,
		AdminSecret: s,
	}, nil
}
