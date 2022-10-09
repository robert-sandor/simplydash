package internal

type Provider interface {
	Get() []Category
}
