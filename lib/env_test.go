package lib_test

import (
	"github.com/mrasu/GravityR/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "test_dummy"
	os.Setenv(key, "dummy")

	v, err := lib.GetEnv(key)
	require.NoError(t, err)
	assert.Equal(t, "dummy", v)
}

func TestGetEnv_error(t *testing.T) {
	v, err := lib.GetEnv("invalid_key")
	assert.ErrorContains(t, err, "not set")
	assert.Equal(t, v, "")
}
