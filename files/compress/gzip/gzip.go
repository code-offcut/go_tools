package gzip

import (
	"archive/tar"
	"fmt"
	"github.com/klauspost/pgzip"
	"github.com/pkg/errors"
	paths2 "go_tools/files/paths"
	"go_tools/log"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type GzipInfo struct {
	SourcePath       string `json:"source_path"`
	TargetPath       string `json:"target_path"`
	Parallelism      int    `json:"parallelism"`
	BlockSize        int64  `json:"block_size"`
	IsCompress       bool   `json:"is_compress"`
	IsDir            bool   `json:"is_dir"`
	IgnoreFailedFile bool   `json:"ignore_failed_file"`
}

func Get(sourcePath string, targetPath string, parallelism int, blockSize int64, isCompress, ignoreFailedFile bool) (result *GzipInfo, err error) {
	if parallelism <= 0 {
		parallelism = runtime.NumCPU()
	}
	if blockSize <= 0 {
		blockSize = 10 * 1024 * 1024
	}
	result = &GzipInfo{
		SourcePath:       "",
		TargetPath:       "",
		Parallelism:      parallelism,
		BlockSize:        blockSize,
		IsCompress:       isCompress,
		IgnoreFailedFile: ignoreFailedFile,
	}
	if len(sourcePath) == 0 {
		return nil, errors.New("source path is empty")
	}
	if sourcePath, err = filepath.Abs(sourcePath); err != nil {
		return nil, errors.Wrap(err, "format source path error")
	} else {
		result.SourcePath = sourcePath
		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, errors.Wrap(err, "source path file/directory not exist")
			} else {
				return nil, errors.Wrap(err, "get source path info error")
			}
		} else {
			result.IsDir = fileInfo.IsDir()
		}
		if !isCompress && result.IsDir {
			return nil, errors.New("source file is directory")
		}
	}
	if len(targetPath) == 0 {
		targetPath = "./"
	}
	if targetPath, err = filepath.Abs(targetPath); err != nil {
		return nil, errors.Wrap(err, "format target path error")
	} else {
		if paths2.IsDir(targetPath) { // target 是目录
			if isCompress {
				if result.IsDir {
					targetPath = filepath.Join(targetPath, "target.tar.gz") // 默认压缩包名称
				} else {
					_, sourceFileName := filepath.Split(result.SourcePath)
					targetPath = filepath.Join(targetPath, fmt.Sprintf("%v.tar.gz", strings.Split(sourceFileName, ".")[0]))
				}
			}
		}
		result.TargetPath = targetPath

	}
	return result, nil
}

func (g *GzipInfo) Compress() error {
	// file write
	fw, err := paths2.CreateFile(g.TargetPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fw.Close(); err != nil {
			log.Debug("close target file error")
		}
	}()
	gzWriter := pgzip.NewWriter(fw)
	defer func() {
		if err := gzWriter.Close(); err != nil {
			log.Debug("close gzip writer error")
		}
	}()
	gzWriter.Comment = "file"
	if g.IsDir {
		gzWriter.Comment = "dir"
	}

	// tar write
	tarWriter := tar.NewWriter(gzWriter)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			log.Debug("close tar writer error")
		}
	}()
	if g.IsDir {
		_, err := paths2.DoIterPath(g.SourcePath, func(fileInfo *paths2.FileInfo, iterErr error) error {
			if iterErr != nil {
				log.Warn(iterErr.Error())
				return nil
			}
			if !fileInfo.IsDir {
				if err := g.gzip(fileInfo, tarWriter); err != nil {
					log.Warn(err.Error())
					if !g.IgnoreFailedFile {
						panic(err)
					}
				}
			}
			return nil
		}, false, g.IgnoreFailedFile)
		return err
	} else {
		sourceFileInfo, err := paths2.GetFileInfo(g.SourcePath)
		if err != nil {
			return err
		}
		if err := g.gzip(sourceFileInfo, tarWriter); err != nil {
			return err
		}
	}
	return nil
}

func (g *GzipInfo) Decompress() error {
	return g.unzip()
}

func (g *GzipInfo) gzip(sourceFile *paths2.FileInfo, writer *tar.Writer) error {
	file, err := os.Open(sourceFile.Path)
	if err != nil {
		return errors.Wrapf(err, "open source file %v error", sourceFile.Path)
	}
	relPath, err := filepath.Rel(g.SourcePath, sourceFile.Path)
	if err != nil {
		return errors.Wrapf(err, "get %v relate path error", sourceFile.Path)
	}
	h := new(tar.Header)
	h.Name = relPath
	h.Size = sourceFile.Size
	h.Mode = int64(os.ModePerm)
	h.ModTime = time.Unix(sourceFile.UpdatedAt, 0)

	if err := writer.WriteHeader(h); err != nil {
		return errors.Wrapf(err, "write file %v header error", sourceFile.Path)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return errors.Wrapf(err, "encode source file %v error", sourceFile.Path)
	}
	return nil
}

func (g *GzipInfo) unzip() error {
	sourceFile, err := os.Open(g.SourcePath)
	if err != nil {
		return errors.Wrapf(err, "open source file %v error", g.SourcePath)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Debug("close source file error: %v", err)
		}
	}()
	gzipReader, err := pgzip.NewReader(sourceFile)
	if err != nil {
		return errors.Wrap(err, "create gzip reader error")
	}
	reader := tar.NewReader(gzipReader)
	fileNumber := 0
	for {
		header, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Warn("read source file error")
				return errors.Wrap(err, "read source file error")
			}
		}
		decodeFilePath := filepath.Join(g.TargetPath, header.Name)
		file, err := paths2.CreateFile(decodeFilePath)
		if err != nil {
			log.Warn(err.Error())
			if !g.IgnoreFailedFile {
				return err
			}
		}
		_, err = io.Copy(file, reader)
		if err != nil {
			log.Warn("writer file %v error: %v", decodeFilePath, err)
			if !g.IgnoreFailedFile {
				return errors.Wrapf(err, "writer file %v error", decodeFilePath)
			}
		}
		fileNumber += 1
	}
	return nil
}
