package internal

type Settings struct {
	Name                   string `json:"name" yaml:"name"`
	EnableHealthIndicators bool   `json:"enable_health_indicators" yaml:"enable_health_indicators"`
	Theme                  Theme  `json:"theme" yaml:"theme"`
}

func DefaultSettings() Settings {
	return Settings{
		Name:                   "simplydash",
		EnableHealthIndicators: false,
		Theme:                  DefaultTheme(),
	}
}

type Theme struct {
	Dark  ThemeColors `json:"dark" yaml:"dark"`
	Light ThemeColors `json:"light" yaml:"light"`
}

func DefaultTheme() Theme {
	return Theme{
		Dark:  DefaultDarkColors(),
		Light: DefaultLightColors(),
	}
}

type ThemeColors struct {
	Background          string `json:"background" yaml:"background"`
	ElementBackground   string `json:"element_background" yaml:"element_background"`
	Foreground          string `json:"foreground" yaml:"foreground"`
	ForegroundSecondary string `json:"foreground_secondary" yaml:"foreground_secondary"`
	AccentColor         string `json:"accent_color" yaml:"accent_color"`
	SuccessColor        string `json:"success_color" yaml:"success_color"`
	WarningColor        string `json:"warning_color" yaml:"warning_color"`
	ErrorColor          string `json:"error_color" yaml:"error_color"`
}

func DefaultDarkColors() ThemeColors {
	return ThemeColors{
		Background:          "rgb(36, 36, 36)",
		ElementBackground:   "rgb(48, 48, 48)",
		Foreground:          "rgba(255, 255, 255, 0.8)",
		ForegroundSecondary: "rgba(255, 255, 255, 0.5)",
		AccentColor:         "rgb(120, 174, 237)",
		SuccessColor:        "rgb(38, 162, 105)",
		WarningColor:        "rgb(205, 147, 9)",
		ErrorColor:          "rgb(192, 28, 40)",
	}
}

func DefaultLightColors() ThemeColors {
	return ThemeColors{
		Background:          "rgb(235, 235, 235)",
		ElementBackground:   "rgb(250, 250, 250)",
		Foreground:          "rgba(0, 0, 0, 0.8)",
		ForegroundSecondary: "rgba(0, 0, 0, 0.5)",
		AccentColor:         "rgb(28, 113, 216)",
		SuccessColor:        "rgb(46, 194, 126)",
		WarningColor:        "rgb(229, 165, 10)",
		ErrorColor:          "rgb(224, 27, 36)",
	}
}
