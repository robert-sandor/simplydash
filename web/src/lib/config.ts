export class Config {
    settings: Settings

    static default(): Config {
        return new Config(Settings.default())
    }

    constructor(settings: Settings) {
        this.settings = settings;
    }
}

export class Settings {
    name: string
    enable_health_indicators: boolean
    theme: Theme

    static default(): Settings {
        return new Settings('simplydash', false, Theme.default())
    }

    constructor(name: string, enable_health_indicators: boolean, theme: Theme) {
        this.name = name;
        this.enable_health_indicators = enable_health_indicators;
        this.theme = theme;
    }
}

export class Theme {
    light: ThemeColors
    dark: ThemeColors

    static default(): Theme {
        return new Theme(
            new ThemeColors(
                "rgb(235, 235, 235)",
                "rgb(250, 250, 250)",
                "rgba(0, 0, 0, 0.8)",
                "rgba(0, 0, 0, 0.5)",
                "rgb(28, 113, 216)",
                "rgb(46, 194, 126)",
                "rgb(229, 165, 10)",
                "rgb(224, 27, 36)",
            ),
            new ThemeColors(
                "rgb(36, 36, 36)",
                "rgb(48, 48, 48)",
                "rgba(255, 255, 255, 0.8)",
                "rgba(255, 255, 255, 0.5)",
                "rgb(120, 174, 237)",
                "rgb(38, 162, 105)",
                "rgb(205, 147, 9)",
                "rgb(192, 28, 40)",
            )
        );
    }

    constructor(light: ThemeColors, dark: ThemeColors) {
        this.light = light;
        this.dark = dark;
    }
}

export class ThemeColors {
    background: string
    element_background: string
    foreground: string
    foreground_secondary: string
    accent_color: string
    success_color: string
    warning_color: string
    error_color: string

    constructor(background: string, element_background: string, foreground: string, foreground_secondary: string, accent_color: string, success_color: string, warning_color: string, error_color: string) {
        this.background = background;
        this.element_background = element_background;
        this.foreground = foreground;
        this.foreground_secondary = foreground_secondary;
        this.accent_color = accent_color;
        this.success_color = success_color;
        this.warning_color = warning_color;
        this.error_color = error_color;
    }
}
