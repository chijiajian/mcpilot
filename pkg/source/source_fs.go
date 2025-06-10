package source

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type FSReader struct {
	SourcePath string
}

func NewFSReader(path string) *FSReader {
	return &FSReader{SourcePath: path}
}

func (f *FSReader) Fetch(ctx context.Context, destDir string) error {
	return copyDir(f.SourcePath, destDir)
}

// 辅助函数：递归复制文件夹
func copyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// 复制文件
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}
