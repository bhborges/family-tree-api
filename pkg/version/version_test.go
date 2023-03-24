package version

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	gitCommit = "f611fad771ce4b11d0f2a579f779faea96b5aff2"
	buildDate = "2023-03-23T01:21:13Z"

	assert.Equal(t, gitCommit, Get().GitCommit)
	assert.Equal(t, buildDate, Get().BuildDate)
	assert.Equal(t, runtime.Version(), Get().GoVersion)
}
