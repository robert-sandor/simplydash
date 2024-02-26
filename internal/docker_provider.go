package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"log/slog"
	"reflect"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	simplydash                    = "simplydash"
	simplydashEnable              = simplydash + ".enable"
	simplydashName                = simplydash + ".name"
	simplydashLink                = simplydash + ".link"
	simplydashGroup               = simplydash + ".group"
	simplydashIcon                = simplydash + ".icon"
	simplydashDescription         = simplydash + ".description"
	simplydashHealthcheckEnable   = simplydash + ".healthcheck.enable"
	simplydashHealthcheckInterval = simplydash + ".healthcheck.interval"
	simplydashHealthcheckTimeout  = simplydash + ".healthcheck.timeout"
)

type DockerProviderConfig struct {
	Host     string        `json:"host" yaml:"host"`
	Interval time.Duration `json:"interval" yaml:"interval"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`
}

type DockerProvider struct {
	clientFunc       DockerClientFunc
	logger           *slog.Logger
	notificationChan chan<- string
	id               string
	apps             []App
	config           DockerProviderConfig
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
		logger:           slog.With("id", id),
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
		<-ticker.C
		dp.fetch()
	}
}

func (dp *DockerProvider) fetch() {
	dockerClient, err := dp.clientFunc(dp.config)
	if err != nil {
		dp.logger.Error("docker client", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dp.config.Interval)
	defer cancel()

	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.NewArgs(filters.Arg("label", simplydashEnable)),
	})
	if err != nil {
		dp.logger.Error("list containers", "error", err)
		return
	}

	apps := make([]App, 0)
	for _, ct := range containers {
		app := dp.containerToApp(ct)
		errs := app.Validate()
		if len(errs) > 0 {
			dp.logger.Error("invalid app specification", "error", errors.Join(errs...))
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
		slog.Error("invalid bool value for label", "error", err, "label", label)
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
		slog.Error("invalid bool value for label", "error", err, "label", label)
		return defaultValue
	}

	if durationVal <= 0 {
		slog.Error("expected positive duration for label", "error", err, "label", label)
		return defaultValue
	}

	return durationVal
}
