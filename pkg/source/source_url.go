package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type URLSource struct {
	URL      string
	Filename string // 可以从 URL 自动提取
}

func NewURLSource(url, filename string) *URLSource {
	return &URLSource{
		URL:      url,
		Filename: filename,
	}
}

func (u *URLSource) Fetch(ctx context.Context, destDir string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", u.URL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, u.URL)
	}

	outFile := filepath.Join(destDir, u.Filename)
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
