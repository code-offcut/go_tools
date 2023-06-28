package copy

import (
	"go_tools/paths"
)

type FailItem struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type TodoFile struct {
	FileInfo   *paths.FileInfo `json:"file_info"`
	TargetPath string          `json:"target_path"`
}
