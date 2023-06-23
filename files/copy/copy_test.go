package copy

import (
	"github.com/stretchr/testify/assert"
	"go_tools/log"
	"testing"
)

func TestCopyV2(t *testing.T) {
	handler, err := Get("/Volumes/tlz/mac/works", "/Volumes/tanlizhi/works_v5", 16)
	assert.Nil(t, err)
	failItems, err := handler.Copy()
	assert.Nil(t, err)
	log.Info("copy process end: success dir number %v, success file number %v, failed number: %v", handler.DirNumber, handler.FileNumber, len(handler.FailFiles))
	if len(failItems) > 0 {
		log.Warn("copy failed files list\n")
		for _, item := range failItems {
			log.Warn("file: %v, reason: %v", item.Path, item.Reason)
		}
	} else {
		log.Info("copy all files success")
	}
}
