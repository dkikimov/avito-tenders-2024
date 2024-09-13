package integration_tests

import (
	"embed"
	"io"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed queries
var Queries embed.FS

type QueriesLoader struct {
	t           *testing.T
	currentPath fs.FS
}

func NewQueriesLoader(t *testing.T, queriesPath fs.FS) *QueriesLoader {
	return &QueriesLoader{
		t:           t,
		currentPath: queriesPath,
	}
}

func (l *QueriesLoader) LoadString(path string) string {
	file, err := l.currentPath.Open(path)
	require.NoError(l.t, err)

	defer file.Close()

	data, err := io.ReadAll(file)
	require.NoError(l.t, err)

	return string(data)
}
