package thelper

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func DecodeHex(t *testing.T, h string) []byte {
	t.Helper()

	res, err := hex.DecodeString(h)
	require.NoError(t, err)
	return res
}
