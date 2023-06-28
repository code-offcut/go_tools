package paths

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIterPath(t *testing.T) {
	files, err := DoIterPath("/Users/tlz/Downloads", nil, true, false)
	assert.Nil(t, err)
	bytes, err := json.Marshal(files)
	assert.Nil(t, err)
	t.Logf("%v", string(bytes))
}
