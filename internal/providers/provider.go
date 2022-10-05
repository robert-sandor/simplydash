package providers

import "simplydash/internal/models"

type Provider interface {
	Load() error
	Get() []models.Category
}
