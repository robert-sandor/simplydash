package internal

import log "github.com/sirupsen/logrus"

type AppService interface {
	Init()
	GetApps() []AppGroup
	UpdateCh() <-chan struct{}
}

func NewAppService(config Config, healthCheckService HealthcheckService) AppService {
	providerUpdateCh := make(chan string, 1)
	providers := BuildProviders(config, providerUpdateCh)

	return &appServiceImpl{
		config:             config,
		healthCheckService: healthCheckService,
		appsByProviderId:   make(map[string][]App),
		providers:          providers,
		providerUpdateCh:   providerUpdateCh,
		updateCh:           make(chan struct{}, 1),
		logger:             log.WithField("name", "app-service"),
	}
}

type appServiceImpl struct {
	config             Config
	appsByProviderId   map[string][]App
	providers          map[string]Provider
	providerUpdateCh   <-chan string
	healthCheckService HealthcheckService
	updateCh           chan struct{}
	logger             *log.Entry
}

func (svc *appServiceImpl) GetApps() []AppGroup {
	indexByGroupName := make(map[string]int)
	appGroups := make([]AppGroup, 0)

	for _, groupName := range svc.config.App.Groups {
		indexByGroupName[groupName] = len(appGroups)
		appGroups = append(appGroups, NewAppGroup(groupName))
	}

	for _, providerApps := range svc.appsByProviderId {
		for _, providerApp := range providerApps {
			var index int
			index, ok := indexByGroupName[providerApp.Group]

			if !ok {
				index = len(appGroups)
				indexByGroupName[providerApp.Group] = index
				appGroups = append(appGroups, NewAppGroup(providerApp.Group))
			}

			if providerApp.Healthcheck.Enabled {
				providerApp.Healthcheck.Health = svc.healthCheckService.Get(providerApp.Link)
			}

			appGroups[index].Apps = insertOrdered(appGroups[index].Apps, providerApp)
		}
	}

	return appGroups
}

func (svc *appServiceImpl) Init() {
	for _, provider := range svc.providers {
		err := provider.Init()
		if err != nil {
			log.WithField("providerId", provider.ID()).WithError(err).Error("init provider")
			continue
		}

		go svc.updateApps(provider.ID())
	}
	go svc.listen()
}

func (svc *appServiceImpl) UpdateCh() <-chan struct{} {
	return svc.updateCh
}

func (svc *appServiceImpl) listen() {
	for {
		select {
		case providerId := <-svc.providerUpdateCh:
			log.WithField("providerId", providerId).Debug("got update from provider")
			go svc.updateApps(providerId)
		case <-svc.healthCheckService.Updates():
			svc.logger.Debug("got update from healthcheck")
			go svc.notify()
		}
	}
}

func (svc *appServiceImpl) updateApps(id string) {
	svc.appsByProviderId[id] = svc.providers[id].Apps()

	svc.refreshHealthCheckers()

	svc.notify()
}

func (svc *appServiceImpl) notify() {
	log.Debug("sending update notification")
	svc.updateCh <- struct{}{}
}

func (svc *appServiceImpl) refreshHealthCheckers() {
	newUrls := make(map[string]AppHealthcheck)
	for _, apps := range svc.appsByProviderId {
		for _, app := range apps {
			if app.Healthcheck.Enabled {
				newUrls[app.Link] = app.Healthcheck
			}
		}
	}

	existingUrls := svc.healthCheckService.Urls()

	for url := range existingUrls {
		if _, ok := newUrls[url]; !ok {
			svc.healthCheckService.Remove(url)
		}
	}

	for url, cfg := range newUrls {
		svc.healthCheckService.Add(url, cfg.Interval, cfg.Timeout)
	}
}
