package models

type Category struct {
	Name  string `json:"name" yaml:"name"`
	Items []Item `json:"items,omitempty" yaml:"items,omitempty"`
}

func NewCategory(name string, items []Item) *Category {
	return &Category{Name: name, Items: items}
}
