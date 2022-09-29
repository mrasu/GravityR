package thelper

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func CreateTemp(t *testing.T, pattern string, fn func(*os.File)) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", pattern)
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	fn(tmpfile)
}
