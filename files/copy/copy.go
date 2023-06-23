package copy

import (
	"fmt"
	"github.com/pkg/errors"
	paths2 "go_tools/files/paths"
	"go_tools/gopool"
	"go_tools/log"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	maxCacheSize = 110 * 1024 * 1024
)

type CopyFiles struct {
	SourcePath  string        `json:"source_path"`
	TargetPath  string        `json:"target_path"`
	Parallelism int           `json:"parallelism"`
	FailFiles   []FailItem    `json:"fail_files"`
	FileNumber  int64         `json:"file_number"`
	DirNumber   int64         `json:"dir_number"`
	todoFileNum int           `json:"-"`
	todoFiles   chan TodoFile `json:"-"`
	endFlag     bool          `json:"-"`
}

func Get(oriPath, targetPath string, parallelism int) (*CopyFiles, error) {
	if len(oriPath) == 0 {
		return nil, errors.New("source path is empty")
	}
	if len(targetPath) == 0 {
		return nil, errors.New("target path is empty")
	}
	if parallelism <= 0 {
		parallelism = runtime.NumCPU() * 2
	}
	return &CopyFiles{
		SourcePath:  oriPath,
		TargetPath:  targetPath,
		Parallelism: parallelism,
		FailFiles:   []FailItem{},
		todoFiles:   make(chan TodoFile, parallelism*2),
		endFlag:     false,
	}, nil
}

func (c *CopyFiles) Copy() ([]FailItem, error) {
	oriPathFormatted, err := filepath.Abs(c.SourcePath)
	if err != nil {
		return nil, errors.Wrapf(err, "format source path %v error", c.SourcePath)
	}
	targetPathFormatted, err := filepath.Abs(c.TargetPath)
	if err != nil {
		return nil, errors.Wrapf(err, "format target path %v error", c.TargetPath)
	}
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for {
			<-ticker.C
			log.Info("copy success dir number: %v && file number: %v, failed number: %v", c.DirNumber, c.FileNumber, len(c.FailFiles))
		}
	}()

	// 消费者，执行复制文件逻辑
	go c.copyFile()

	// 生产者，生产需要 copy 的文件列表
	_, err = paths2.DoIterPath(c.SourcePath, func(fileInfo *paths2.FileInfo, iterErr error) error {
		if iterErr != nil {
			c.FailFiles = append(c.FailFiles, FailItem{
				Path:   fileInfo.Path,
				Reason: iterErr.Error(),
			})
			return nil
		}
		defer func() {
			if err := recover(); err != nil {
				c.FailFiles = append(c.FailFiles, FailItem{
					Path:   fileInfo.Path,
					Reason: fmt.Sprintf("%v", err),
				})
				log.Warn("%v", err)
			}
		}()
		targetFilePath := strings.ReplaceAll(fileInfo.Path, oriPathFormatted, targetPathFormatted)
		if err = c.doCopy(fileInfo, targetFilePath); err != nil {
			c.FailFiles = append(c.FailFiles, FailItem{
				Path:   fileInfo.Path,
				Reason: err.Error(),
			})
			return nil
		}
		return nil
	}, false, true)
	if err != nil {
		return nil, err
	}
	c.endFlag = true
	return c.FailFiles, nil
}

func (c *CopyFiles) doCopy(fileInfo *paths2.FileInfo, targetPath string) error {
	if fileInfo.IsDir {
		c.DirNumber += 1
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return errors.Wrapf(err, "make dir %v error", targetPath)
		}
		if len(fileInfo.Includes) > 0 {
			for _, info := range fileInfo.Includes {
				if err := c.doCopy(info, filepath.Join(targetPath, info.Name)); err != nil {
					return err
				}
			}
		}
	} else {
		c.FileNumber += 1
		c.todoFileNum += 1
		c.todoFiles <- TodoFile{
			FileInfo:   fileInfo,
			TargetPath: targetPath,
		}
	}
	return nil
}

func (c *CopyFiles) copyFile() {
	pool := gopool.New(c.Parallelism)
	for true {
		fileInfo := <-c.todoFiles
		c.todoFileNum -= 1
		go func() {
			pool.Add(1)
			defer pool.Done()
			sourceFile, err := os.Open(fileInfo.FileInfo.Path)
			if err != nil {
				c.FailFiles = append(c.FailFiles, FailItem{
					Path:   fileInfo.FileInfo.Path,
					Reason: fmt.Sprintf("open source file error: %v ", err),
				})
				return
			}
			defer func() {
				if err := sourceFile.Close(); err != nil {
					log.Debug(err.Error())
				}
			}()
			destination, err := os.Create(fileInfo.TargetPath)
			if err != nil {
				c.FailFiles = append(c.FailFiles, FailItem{
					Path:   fileInfo.FileInfo.Path,
					Reason: fmt.Sprintf("open target file path %v error: %v", fileInfo.TargetPath, err),
				})
				return
			}
			defer func() {
				if err := destination.Close(); err != nil {
					log.Debug(err.Error())
				}
			}()
			var cache []byte
			if maxCacheSize > fileInfo.FileInfo.Size {
				cache = make([]byte, fileInfo.FileInfo.Size+1)
			} else {
				cache = make([]byte, maxCacheSize)
			}
			_, err = io.CopyBuffer(destination, sourceFile, cache)
			if err != nil {
				c.FailFiles = append(c.FailFiles, FailItem{
					Path:   fileInfo.FileInfo.Path,
					Reason: fmt.Sprintf("copy file error: %v", err),
				})
				return
			}
		}()
		if c.endFlag && c.todoFileNum == 0 { // 遍历已经结束，而且所有待 copy 文件已经进入 copy 流程。只需要等待所有 copy 协程结束就可以了
			pool.Wait()
			return
		}
	}
}
