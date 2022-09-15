package thelper

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"
)

func ReadFromFiles(t *testing.T, filename string) string {
	t.Helper()

	_, b, _, _ := runtime.Caller(0)
	p := filepath.Join(filepath.Dir(b), "../", "./testdata/files/", filename)
	txt, err := ioutil.ReadFile(p)
	require.NoError(t, err)

	return string(txt)
}
