package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"time"
)

const simplydash = "simplydash"
const simplydashEnable = simplydash + ".enable"
const simplydashName = simplydash + ".name"
const simplydashLink = simplydash + ".link"
const simplydashGroup = simplydash + ".group"
const simplydashIcon = simplydash + ".icon"
const simplydashDescription = simplydash + ".description"
const simplydashHealthcheckEnable = simplydash + ".healthcheck.enable"
const simplydashHealthcheckInterval = simplydash + ".healthcheck.interval"
const simplydashHealthcheckTimeout = simplydash + ".healthcheck.timeout"

type DockerProviderConfig struct {
	Host     string        `json:"host" yaml:"host"`
	Interval time.Duration `json:"interval" yaml:"interval"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`
}

type DockerProvider struct {
	id               string
	apps             []App
	config           DockerProviderConfig
	clientFunc       DockerClientFunc
	logger           *log.Entry
	notificationChan chan<- string
}

type DockerClientFunc func(config DockerProviderConfig) (client.APIClient, error)

func RealDockerClientFunc() DockerClientFunc {
	return func(config DockerProviderConfig) (client.APIClient, error) {
		return client.NewClientWithOpts(
			client.WithHost(config.Host),
			client.WithTimeout(config.Timeout),
			client.WithAPIVersionNegotiation(),
		)
	}
}

func NewDockerProvider(name string, config DockerProviderConfig, notificationChan chan<- string) Provider {
	id := fmt.Sprintf("docker-%s", name)
	return &DockerProvider{
		id:               id,
		apps:             make([]App, 0),
		config:           config,
		clientFunc:       RealDockerClientFunc(),
		notificationChan: notificationChan,
		logger:           log.WithField("id", id),
	}
}

func (dp *DockerProvider) ID() string {
	return dp.id
}

func (dp *DockerProvider) Apps() []App {
	return dp.apps
}

func (dp *DockerProvider) Init() error {
	go dp.fetch()
	go dp.poll()
	return nil
}

func (dp *DockerProvider) poll() {
	if dp.config.Interval <= 0 {
		dp.config.Interval = time.Minute
	}

	ticker := time.NewTicker(dp.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dp.fetch()
		}
	}
}

func (dp *DockerProvider) fetch() {
	dockerClient, err := dp.clientFunc(dp.config)
	if err != nil {
		dp.logger.WithError(err).Error("docker client")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dp.config.Interval)
	defer cancel()

	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.NewArgs(filters.Arg("label", simplydashEnable)),
	})
	if err != nil {
		dp.logger.WithError(err).Error("list containers")
		return
	}

	apps := make([]App, 0)
	for _, container := range containers {
		app := dp.containerToApp(container)
		errs := app.Validate()
		if len(errs) > 0 {
			dp.logger.WithError(errors.Join(errs...)).Error("invalid app specification")
		} else {
			apps = insertOrdered(apps, app)
		}
	}

	if !reflect.DeepEqual(dp.apps, apps) {
		dp.apps = apps
		dp.notificationChan <- dp.id
	}
}

func (dp *DockerProvider) containerToApp(container types.Container) App {
	return App{
		Name:        container.Labels[simplydashName],
		Description: container.Labels[simplydashDescription],
		Link:        container.Labels[simplydashLink],
		Icon:        container.Labels[simplydashIcon],
		Group:       container.Labels[simplydashGroup],
		Healthcheck: AppHealthcheck{
			Enabled:  boolFromLabel(container, simplydashHealthcheckEnable, DefaultEnableHealthcheck),
			Health:   Unknown,
			Interval: durationFromLabel(container, simplydashHealthcheckInterval, DefaultHealthcheckInterval),
			Timeout:  durationFromLabel(container, simplydashHealthcheckTimeout, DefaultHealthcheckTimeout),
		},
	}
}

func boolFromLabel(container types.Container, label string, defaultValue bool) bool {
	stringVal, ok := container.Labels[label]
	if !ok {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(stringVal)
	if err != nil {
		log.WithError(err).Warnf("invalid bool value for label '%s'", label)
		return defaultValue
	}
	return boolVal
}

func durationFromLabel(container types.Container, label string, defaultValue time.Duration) time.Duration {
	stringVal, ok := container.Labels[label]
	if !ok {
		return defaultValue
	}

	durationVal, err := time.ParseDuration(stringVal)
	if err != nil {
		log.WithError(err).Warnf("invalid bool value for label '%s'", label)
		return defaultValue
	}

	if durationVal <= 0 {
		log.Warnf("expected positive duration for label '%s'", label)
		return defaultValue
	}

	return durationVal
}
