package bootstrap

import (
	"context"
	"search-country-system/internal/app/handlers"
	"search-country-system/internal/app/services"
)

type App struct {
	CountryHandler *handlers.CountryHandler
}

func NewApp(ctx context.Context) *App {
	// dependencies
	cache := services.InitCacheService()
	serviceCall := services.InitServiceCall()
	countryService := services.InitCountryService(serviceCall, cache)

	// handlers
	countryHandler := handlers.InitCountryHandler(countryService)

	return &App{
		CountryHandler: countryHandler,
	}
}
