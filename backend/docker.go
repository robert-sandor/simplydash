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

type DockerProvider struct {
	apps       []App
	configs    []DockerConfig
	stopSignal chan struct{}
}

type fetchResult struct {
	apps []App
	err  []error
}

func NewDockerProvider(configs []DockerConfig) *DockerProvider {
	return &DockerProvider{
		apps:       make([]App, 0),
		configs:    configs,
		stopSignal: make(chan struct{}, 1),
	}
}

func Init() error {
	return nil
}

func (dockerProvider *DockerProvider) Stop() {
	dockerProvider.stopSignal <- struct{}{}
}

func (dockerProvider *DockerProvider) fetchApps() fetchResult {
	waitGroup := sync.WaitGroup{}

	resultByConfig := make(map[*DockerConfig]chan fetchResult)
	for _, config := range dockerProvider.configs {
		resultByConfig[&config] = make(chan fetchResult, 1)
		waitGroup.Add(1)
		go fetch(&config, resultByConfig[&config], &waitGroup)
	}

	return fetchResult{}
}

func fetch(dockerConfig *DockerConfig, result chan fetchResult, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	dockerClient, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
		client.WithHost(dockerConfig.Host),
	)
	if err != nil {
		result <- fetchResult{
			apps: []App{},
			err:  []error{err},
		}
	}

	containers, err := func() ([]types.Container, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return dockerClient.ContainerList(ctx, types.ContainerListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("label", enabledLabel)),
		})
	}()
	if err != nil {
		result <- fetchResult{
			apps: []App{},
			err:  []error{err},
		}
	}

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
