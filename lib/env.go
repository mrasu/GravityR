package lib

import (
	"github.com/pkg/errors"
	"os"
)

func GetEnv(key string) (string, error) {
	v := os.Getenv(key)
	if len(v) == 0 {
		return "", errors.Errorf("env: %s is not set", key)
	}
	return v, nil
}
