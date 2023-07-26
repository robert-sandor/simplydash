package internal

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
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
	err := os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	response, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer closeSafe(response.Body)

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
		log.WithError(err).Warn("closing resource")
	}
}
