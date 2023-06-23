package gzip

import (
	"archive/tar"
	"compress/gzip"
	"github.com/klauspost/pgzip"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestGzip(t *testing.T) {
	// file write
	fw, err := os.Create("golang_src.tar.gz")
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// open file to compress
	fileInfo, err := os.Stat("~/logs.txt")
	assert.Nil(t, err)
	file, err := os.Open("~/logs.txt")
	assert.Nil(t, err)
	defer file.Close()
	h := new(tar.Header)
	h.Name = fileInfo.Name()
	h.Size = fileInfo.Size()
	h.Mode = int64(fileInfo.Mode())
	h.ModTime = fileInfo.ModTime()

	err = tw.WriteHeader(h)
	assert.Nil(t, err)

	_, err = io.Copy(tw, file)
	assert.Nil(t, err)

	t.Logf("compress success")

}

func TestPGzip(t *testing.T) {
	// file write
	fw, err := os.Create("sgolang_src.tar.gz")
	if err != nil {
		panic(err)
	}
	defer fw.Close()
	writer := pgzip.NewWriter(fw)
	defer writer.Close()

	// tar write
	tw := tar.NewWriter(writer)
	defer tw.Close()

	// open file to compress
	fileInfo, err := os.Stat("~/logs.txt")
	assert.Nil(t, err)
	file, err := os.Open("~/logs.txt")
	assert.Nil(t, err)
	defer file.Close()
	h := new(tar.Header)
	h.Name = fileInfo.Name()
	h.Size = fileInfo.Size()
	h.Mode = int64(os.ModePerm)
	h.ModTime = fileInfo.ModTime()

	err = tw.WriteHeader(h)
	assert.Nil(t, err)

	_, err = io.Copy(tw, file)
	assert.Nil(t, err)

}

func TestCompressDir(t *testing.T) {
	handler, err := Get("~/Downloads/", "~/test.tar.gz", 0, 0, true, false)
	assert.Nil(t, err)
	err = handler.Compress()
	assert.Nil(t, err)
}

func TestCompress(t *testing.T) {
	handler, err := Get("~/Downloads/0000.txt", "~/test_file.tar.gz", 0, 0, true, false)
	assert.Nil(t, err)
	err = handler.Compress()
	assert.Nil(t, err)
}

func TestDecompressDir(t *testing.T) {
	handler, err := Get("~/test.tar.gz", "~/temp.tar.gz", 0, 0, false, false)
	assert.Nil(t, err)
	err = handler.Decompress()
	assert.Nil(t, err)
}

func TestDecompressFile(t *testing.T) {
	handler, err := Get("~/test_file.tar.gz", "~/temp", 0, 0, false, false)
	assert.Nil(t, err)
	err = handler.Decompress()
	assert.Nil(t, err)
}
