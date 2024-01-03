package file

import (
	"bufio"
	"context"
	"github.com/google/uuid"
	"io"
	"os"
	gopath "path"
	"time"
)

type Storage interface {
	save(ctx context.Context, filename string, reader io.Reader) (string, error)
	delete(ctx context.Context, path string) error
}

type fileStorage struct {
	rootDir string
}

func (fs *fileStorage) save(ctx context.Context, filename string, reader io.Reader) (string, int64, error) {
	uniqueID := uuid.New().String()
	folder := time.Now().Format("20060102/15/04")
	userPath := gopath.Join("/", folder, uniqueID)
	folder = gopath.Join(fs.rootDir, folder)
	err := os.MkdirAll(folder, 0775)
	if err != nil {
		return "", 0, err
	}
	var absPath string
	if filename == "" {
		absPath = gopath.Join(folder, uniqueID)
	} else {
		absPath = gopath.Join(folder, uniqueID)
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", 0, err
	}
	size, err := io.Copy(file, reader)
	if err != nil {
		return "", 0, err
	}
	return userPath, size, nil
}

func (fs *fileStorage) read(ctx context.Context, path string) (io.Reader, error) {
	absPath := gopath.Join(fs.rootDir, path)
	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return reader, nil
}

func (fs *fileStorage) delete(ctx context.Context, path string) error {
	absPath := gopath.Join(fs.rootDir, path)
	err := os.Remove(absPath)
	if err != nil {
		return err
	}
	return nil
}
