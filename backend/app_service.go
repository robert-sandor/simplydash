package main

type AppService interface {
	GetApps() []App
}

type AppServiceImpl struct {
	configService  ConfigService
	dockerProvider DockerProvider
}

func NewAppService(configService ConfigService, dockerProvider DockerProvider) AppService {
	return &AppServiceImpl{
		configService:  configService,
		dockerProvider: dockerProvider,
	}
}

// GetApps implements AppService.
func (as *AppServiceImpl) GetApps() []App {
	dockerApps := as.dockerProvider.GetApps()
	configApps := as.configService.Get().Apps
	return append(dockerApps, configApps...)
}
