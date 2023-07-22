package internal

import "testing"

func TestApp_resolveIconUrl(t *testing.T) {
	tests := []struct {
		name     string
		icon     string
		expected string
	}{
		{
			name:     "simple case",
			icon:     "google",
			expected: "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/google.svg",
		},
		{
			name:     "with whitespace",
			icon:     "adguard home",
			expected: "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/adguard-home.svg",
		},
		{
			name:     "with multiple whitespaces",
			icon:     " adguard _- home ",
			expected: "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/adguard-home.svg",
		},
		{
			name:     "with underscore",
			icon:     "adguard_home",
			expected: "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/adguard-home.svg",
		},
		{
			name:     "with dash",
			icon:     "adguard-home",
			expected: "https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/svg/adguard-home.svg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{Icon: tt.icon}
			app.resolveIconUrl()
			if app.Icon != tt.expected {
				t.Errorf("expected '%s' but got '%s'", tt.expected, app.Icon)
			}
		})
	}
}
