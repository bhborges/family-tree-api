package version

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	gitCommit = "d240853866f20fc3e536cb3bca86c86c54b723ce"
	buildDate = "2021-09-17T01:21:13Z"

	assert.Equal(t, gitCommit, Get().GitCommit)
	assert.Equal(t, buildDate, Get().BuildDate)
	assert.Equal(t, runtime.Version(), Get().GoVersion)
}
