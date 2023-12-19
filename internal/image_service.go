package internal

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
)

type ImageService interface {
	Get(urlString string) (string, error)
}

func NewImageService(cachePath string) ImageService {
	return &imageServiceImpl{
		cachePath: cachePath,
	}
}

type imageServiceImpl struct {
	cachePath string
}

func (svc *imageServiceImpl) Get(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	filePath := path.Join(svc.cachePath, u.Hostname(), u.Path)
	_, err = os.Stat(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		err = svc.downloadImage(u, filePath)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

func (svc *imageServiceImpl) downloadImage(u *url.URL, filePath string) error {
	err := os.MkdirAll(path.Dir(filePath), 0o755)
	if err != nil {
		return err
	}

	response, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer closeSafe(response.Body)

	if response.StatusCode >= 400 {
		return fmt.Errorf("got status code %d", response.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer closeSafe(file)

	_, err = io.Copy(file, response.Body)
	return err
}

func closeSafe(c io.Closer) {
	err := c.Close()
	if err != nil {
		slog.Error("unexpected error when closing resource", "error", err)
	}
}
