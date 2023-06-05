package main

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

const (
	enabledLabel                          = "simplydash.enabled=true"
	nameLabel                             = "simplydash.name"
	descriptionLabel                      = "simplydash.description"
	groupLabel                            = "simplydash.group"
	linkLabel                             = "simplydash.link"
	iconLabel                             = "simplydash.icon"
	healthcheckDisableLabel               = "simplydash.healthcheck.disable"
	healthcheckLinkLabel                  = "simplydash.healthcheck.link"
	healthcheckPollIntervalInSecondsLabel = "simplydash.healthcheck.pollIntervalInSeconds"
)

type DockerProvider interface {
	Init() error
	Stop()
	GetApps() []App
}

type DockerProviderImpl struct {
	apps       []App
	configs    map[string]DockerConfig
	stopSignal chan struct{}
}

type fetchResult struct {
	apps []App
	err  []error
}

func NewDockerProvider(configs map[string]DockerConfig) DockerProvider {
	return &DockerProviderImpl{
		apps:       make([]App, 0),
		configs:    configs,
		stopSignal: make(chan struct{}, 1),
	}
}

func (dp *DockerProviderImpl) GetApps() []App {
	return dp.apps
}

func (dp *DockerProviderImpl) Init() error {
	logrus.Debug("initializing docker provider")
	result := dp.fetchApps()
	if len(result.err) > 0 {
		return errors.Join(result.err...)
	}
	dp.apps = result.apps
	return nil
}

func (dockerProvider *DockerProviderImpl) Stop() {
	dockerProvider.stopSignal <- struct{}{}
}

func (dockerProvider *DockerProviderImpl) fetchApps() fetchResult {
	logrus.Debug("fetching apps from docker")
	waitGroup := sync.WaitGroup{}

	resultByConfig := make(map[string]chan fetchResult)
	for key, config := range dockerProvider.configs {
		resultByConfig[key] = make(chan fetchResult, 1)
		waitGroup.Add(1)
		go fetch(&config, resultByConfig[key], &waitGroup)
	}

	logrus.Debug("waiting for goroutines to finish")
	waitGroup.Wait()

	result := fetchResult{
		apps: make([]App, 0),
		err:  make([]error, 0),
	}
	for _, resultByConfigChan := range resultByConfig {
		select {
		case res, ok := <-resultByConfigChan:
			if !ok {
				result.err = append(result.err, errors.New("no result found"))
				continue
			}

			result.apps = append(result.apps, res.apps...)
			result.err = append(result.err, res.err...)
		}
	}

	return result
}

func fetch(dockerConfig *DockerConfig, result chan fetchResult, waitGroup *sync.WaitGroup) {
	logrus.WithField("config", dockerConfig).Debug("fetching apps from docker")
	defer waitGroup.Done()

	dockerClient, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
		client.WithHost(dockerConfig.Host),
	)
	if err != nil {
		logrus.WithField("config", dockerConfig).WithField("err", err).Error("conneting to docker")
		result <- fetchResult{
			apps: []App{},
			err:  []error{err},
		}
	}
	logrus.WithField("config", dockerConfig).Debug("connected to docker")

	containers, err := func() ([]types.Container, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return dockerClient.ContainerList(ctx, types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("label", enabledLabel)),
		})
	}()
	if err != nil {
		logrus.WithField("config", dockerConfig).WithField("err", err).Error("listing containers")
		result <- fetchResult{
			apps: []App{},
			err:  []error{err},
		}
	}
	logrus.WithField("container_count", len(containers)).Debug("listing containers")

	res := fetchResult{
		apps: make([]App, 0),
		err:  make([]error, 0),
	}
	for _, container := range containers {
		app, err := appFromContainer(container)
		if err != nil {
			res.err = append(res.err, err)
		} else {
			res.apps = append(res.apps, app)
		}
	}

	logrus.WithField("result", res).Debug("returning apps")
	result <- res
}

func appFromContainer(container types.Container) (App, error) {
	app := App{
		Healthcheck: AppHealthcheck{
			PollIntervalInSeconds: 60,
		},
	}

	if name, ok := container.Labels[nameLabel]; ok {
		app.Name = name
	} else {
		return app, errors.New("name is required")
	}

	if link, ok := container.Labels[linkLabel]; ok {
		app.Link = link
	} else {
		return app, errors.New("link is required")
	}

	if group, ok := container.Labels[groupLabel]; ok {
		app.Group = group
	} else {
		return app, errors.New("group is required")
	}

	if description, ok := container.Labels[descriptionLabel]; ok {
		app.Description = description
	} else {
		app.Description = app.Link
	}

	if icon, ok := container.Labels[iconLabel]; ok {
		app.Icon = icon
	} else {
		app.Icon = defaultIcon(app.Name)
	}

	if healthCheckDisable, ok := container.Labels[healthcheckDisableLabel]; ok {
		disable, err := strconv.ParseBool(healthCheckDisable)
		if err == nil {
			app.Healthcheck.Disable = disable
		}
	}

	if healthcheckLink, ok := container.Labels[healthcheckLinkLabel]; ok {
		app.Healthcheck.Link = healthcheckLink
	} else {
		app.Healthcheck.Link = app.Link
	}

	if healthcheckPoll, ok := container.Labels[healthcheckPollIntervalInSecondsLabel]; ok {
		seconds, err := strconv.Atoi(healthcheckPoll)
		if err == nil && seconds > 0 {
			app.Healthcheck.PollIntervalInSeconds = seconds
		}
	}

	return app, nil
}
