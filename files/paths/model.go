package paths

type FileInfo struct {
	Name      string               `json:"name"`
	Path      string               `json:"path"`
	Size      int64                `json:"size"`
	IsDir     bool                 `json:"is_dir"`
	UpdatedAt int64                `json:"updated_at"`
	Parent    *FileInfo            `json:"-"`
	Includes  map[string]*FileInfo `json:"includes"`
}

type DoFunc func(fileInfo *FileInfo, iterErr error) error
