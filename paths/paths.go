package paths

import (
	"fmt"
	"github.com/pkg/errors"
	"go_tools/log"
	"os"
	"path/filepath"
	"sort"
)

func DoIterPath(path string, doFunc DoFunc, includeItems, ignoreErr bool) (result *FileInfo, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); err != nil {
		return nil, errors.Wrapf(err, "get path %v info error", path)
	}
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return nil, errors.Wrapf(err, "get file %v info error", path)
	}
	result = &FileInfo{
		Name:      fileInfo.Name(),
		Path:      path,
		Size:      fileInfo.Size(),
		IsDir:     fileInfo.IsDir(),
		UpdatedAt: fileInfo.ModTime().Unix(),
		Parent:    nil,
		Includes:  map[string]*FileInfo{},
	}

	if doFunc != nil {
		if err = doFunc(result, nil); err != nil {
			return nil, err
		}
	}
	if fileInfo.IsDir() {
		files, err := readDirNames(path)
		if err != nil {
			return nil, errors.Wrap(err, "list sub files error")
		}
		if len(files) > 0 {
			for i := range files {
				walk(result, fmt.Sprintf("%v/%v", path, files[i]), 1, doFunc, includeItems, ignoreErr)
			}
		}
	}
	return result, nil
}

func walk(parent *FileInfo, path string, depth int, doFunc DoFunc, includeItems, ignoreErr bool) {
	defer func() {
		if ignoreErr {
			if err := recover(); err != nil {
				_ = doFunc(parent, fmt.Errorf("iter path %v error: %v", path, err))
				//log.Warn("iter path %v error: %v", path, err)
			}
		}
	}()
	fileInfo, err := os.Lstat(path)
	if err != nil {
		if err := doFunc(&FileInfo{
			Path: path,
		}, fmt.Errorf("get file %v info error: %v", path, err)); err != nil {
			panic(err)
		}
	} else {
		var result = FileInfo{
			Name:      fileInfo.Name(),
			Path:      path,
			Size:      fileInfo.Size(),
			IsDir:     fileInfo.IsDir(),
			UpdatedAt: fileInfo.ModTime().Unix(),
			Parent:    parent,
			Includes:  map[string]*FileInfo{},
		}
		if err = doFunc(&result, nil); err != nil {
			panic(fmt.Sprintf("do func error: %v", err))
		}
		if fileInfo.IsDir() {
			names, err := readDirNames(path)
			if err != nil {
				if err := doFunc(&result, errors.Wrapf(err, "list %v sub files error", path)); err != nil {
					panic(err)
				}
			}
			for i := range names {
				var iterName = names[i]
				walk(&result, fmt.Sprintf("%v/%v", path, iterName), depth+1, doFunc, includeItems, ignoreErr)
			}
		}
		if includeItems {
			parent.Includes[fileInfo.Name()] = &result
		}
	}
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Debug("close file %v error: %v", dirname, err)
		}
	}(f)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)

	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func GetFileInfo(path string) (*FileInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "fotmat path %v error", path)
	}
	fileInfo, err := os.Lstat(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "get file %v info error", absPath)
	}
	result := &FileInfo{
		Name:      fileInfo.Name(),
		Path:      absPath,
		Size:      fileInfo.Size(),
		IsDir:     fileInfo.IsDir(),
		UpdatedAt: fileInfo.ModTime().Unix(),
		Parent:    nil,
		Includes:  map[string]*FileInfo{},
	}
	return result, nil
}

func CreateFile(path string) (*os.File, error) {
	var dirPath string
	if IsDir(path) {
		dirPath = path
	} else {
		dirPath, _ = filepath.Split(path)
	}
	_, err := os.Lstat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return nil, errors.Wrapf(err, "create directory %v error", err)
			}
		} else {
			return nil, errors.Wrapf(err, "get dir %v info error", dirPath)
		}
	}
	if !IsDir(path) {
		file, err := os.Create(path)
		if err != nil {
			return nil, errors.Wrapf(err, "create file %v error", path)
		}
		return file, nil
	}
	return nil, nil
}

func IsDir(path string) bool {
	volumeName := filepath.VolumeName(path)
	dir, _ := filepath.Split(path)
	return len(volumeName) == len(dir)
}
