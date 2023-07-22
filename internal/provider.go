package internal

type Provider interface {
	ID() string
	Apps() []App
	Init() error
}

func BuildProviders(config Config, notificationChan chan<- string) map[string]Provider {
	providers := make(map[string]Provider)

	for providerName, providerConfig := range config.Providers.Docker {
		provider := NewDockerProvider(providerName, providerConfig, notificationChan)
		providers[provider.ID()] = provider
	}

	for providerName, providerConfig := range config.Providers.File {
		provider := NewFileProvider(providerName, providerConfig, notificationChan)
		providers[provider.ID()] = provider
	}

	return providers
}
