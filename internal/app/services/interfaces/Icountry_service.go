package interfaces

import (
	"context"
	"search-country-system/internal/dto"
)

type ICountryService interface {
	GetCountryData(ctx context.Context, name string) (dto.Country, error)
}
