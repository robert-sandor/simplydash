package internal

import (
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type HealthcheckService interface {
	Init()
	Get(url string) AppHealth
	Urls() map[string]bool
	Add(url string, interval time.Duration, timeout time.Duration)
	Remove(url string)
	Updates() <-chan struct{}
}

func NewHealthcheckService() HealthcheckService {
	return &healthcheckServiceImpl{
		checkersByUrl:    make(map[string]*healthChecker),
		checkersUpdateCh: make(chan string, 1),
		updateCh:         make(chan struct{}, 1),
	}
}

type healthcheckServiceImpl struct {
	checkersByUrl    map[string]*healthChecker
	checkersUpdateCh chan string
	updateCh         chan struct{}
}

func (svc *healthcheckServiceImpl) Init() {
	go svc.listen()
}

func (svc *healthcheckServiceImpl) Get(url string) AppHealth {
	if checker, ok := svc.checkersByUrl[url]; ok {
		return checker.health
	}
	return Error
}

func (svc *healthcheckServiceImpl) Urls() map[string]bool {
	urls := make(map[string]bool)
	for url := range svc.checkersByUrl {
		urls[url] = true
	}
	return urls
}

func (svc *healthcheckServiceImpl) Add(url string, interval time.Duration, timeout time.Duration) {
	if checker, ok := svc.checkersByUrl[url]; ok {
		checker.update(interval, timeout)
		return
	}
	svc.checkersByUrl[url] = newHealthChecker(url, interval, timeout, svc.checkersUpdateCh)
}

func (svc *healthcheckServiceImpl) Remove(url string) {
	if checker, ok := svc.checkersByUrl[url]; ok {
		checker.stopCh <- struct{}{}
		delete(svc.checkersByUrl, url)
		return
	}
}

func (svc *healthcheckServiceImpl) Updates() <-chan struct{} {
	return svc.updateCh
}

func (svc *healthcheckServiceImpl) listen() {
	for {
		<-svc.checkersUpdateCh
		go func() {
			svc.updateCh <- struct{}{}
		}()
	}
}

type healthChecker struct {
	updateCh chan<- string
	stopCh   chan struct{}
	ticker   *time.Ticker
	url      string
	interval time.Duration
	timeout  time.Duration
	health   AppHealth
}

func newHealthChecker(
	url string,
	interval time.Duration,
	timeout time.Duration,
	checkerUpdateCh chan<- string,
) *healthChecker {
	checker := &healthChecker{
		url:      url,
		interval: interval,
		timeout:  timeout,
		health:   Unknown,
		updateCh: checkerUpdateCh,
		stopCh:   make(chan struct{}, 1),
	}
	go checker.updateHealth()
	go checker.run()
	return checker
}

func (h *healthChecker) run() {
	h.ticker = time.NewTicker(h.interval)
	defer h.ticker.Stop()

	for {
		select {
		case <-h.stopCh:
			return
		case <-h.ticker.C:
			go h.updateHealth()
		}
	}
}

func (h *healthChecker) update(interval time.Duration, timeout time.Duration) {
	h.ticker.Reset(interval)
	h.timeout = timeout
}

func (h *healthChecker) updateHealth() {
	health := h.healthCheck()
	if h.health != health {
		h.health = health
		h.updateCh <- h.url
	}
}

func (h *healthChecker) healthCheck() AppHealth {
	client := http.Client{Timeout: h.timeout}
	response, err := client.Get(h.url)
	if err != nil {
		log.WithField("url", h.url).WithError(err).Warn("healthcheck error")
		if os.IsTimeout(err) {
			return Timeout
		} else {
			return Error
		}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode >= 400 {
		log.WithField("url", h.url).
			WithField("code", response.StatusCode).
			Warn("healthcheck status")
		if response.StatusCode == 404 || response.StatusCode >= 500 {
			return Error
		}
		return Warning
	}

	log.WithField("url", h.url).Debug("healthcheck")
	return Healthy
}
