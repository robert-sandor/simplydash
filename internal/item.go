package internal

type Item struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func NewItem(name string, URL string, icon string, description string) *Item {
	return &Item{Name: name, URL: URL, Icon: icon, Description: description}
}
