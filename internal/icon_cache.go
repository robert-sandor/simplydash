package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

type IconCache struct {
	path string
}

func NewIconCache(path string) *IconCache {
	return &IconCache{path: path}
}

func (c *IconCache) GetIcon(urlStr string) (string, error) {
	parsedUrl, err := url.ParseRequestURI(urlStr)
	if err != nil {
		iconPath := urlStr
		if !path.IsAbs(iconPath) {
			iconPath = path.Join(c.path, iconPath)
		}

		if _, err := os.Stat(iconPath); err != nil {
			return "", err
		}

		return iconPath, nil
	}

	iconPath := path.Join(c.path, parsedUrl.Host, parsedUrl.Path)

	if _, err := os.Stat(iconPath); err != nil {
		return fetchIcon(parsedUrl, iconPath)
	}

	return iconPath, nil
}

func fetchIcon(parsedUrl *url.URL, iconPath string) (string, error) {
	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body err = %+v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	err = os.MkdirAll(path.Dir(iconPath), 0700)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(iconPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file %s err = %+v", file.Name(), err)
		}
	}(file)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return iconPath, nil
}
