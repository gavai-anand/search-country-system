package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"search-country-system/internal/app/constants"
	"search-country-system/internal/app/services/interfaces"
	"search-country-system/internal/dto"
)

// CountryService provides methods to fetch country data from cache or external REST API.
type CountryService struct {
	sc                interfaces.IServiceCall  // HTTP client interface to make service calls
	cacheService      interfaces.ICacheService // Cache service to store/retrieve country data
	restCountryDomain string                   // Base URL of the REST Countries API
}

// InitCountryService initializes and returns a new CountryService instance
func InitCountryService(sc interfaces.IServiceCall, cache interfaces.ICacheService) *CountryService {
	return &CountryService{
		sc:                sc,
		cacheService:      cache,
		restCountryDomain: viper.GetString("REST_COUNTRIES_SERVICE"), // Read API domain from environment
	}
}

// GetCountryData fetches country data by name.
// It first checks the cache, and if not present, calls the REST API.
func (cs CountryService) GetCountryData(ctx context.Context, name string) (dto.Country, error) {
	var country dto.Country

	// ------------------------
	// Step 1: Check cache
	// ------------------------
	existingData, isOk := cs.cacheService.Get(name)
	if isOk {
		log.Println("Cache hit for country:", name)
		if cachedCountry, ok := existingData.(dto.Country); ok {
			return cachedCountry, nil
		}

	} else {
		log.Println("Cache miss for country:", name)
	}

	// ------------------------
	// Step 2: Call external REST API
	// ------------------------
	restCountryRsp, err := cs.GetCountryDataServiceCall(ctx, name)
	if err != nil {
		log.Printf("Error fetching country data from API for '%s': %v\n", name, err)
		return country, err
	}

	if len(restCountryRsp) == 0 {
		errMsg := fmt.Sprintf("No country data found for '%s'", name)
		log.Println(errMsg)
		return country, fmt.Errorf("%s", errMsg)
	}

	// ------------------------
	// Step 3: Map API response to internal DTO
	// ------------------------
	singleCountry := restCountryRsp[0]

	// Capital: pick the first element if available
	capital := ""
	if len(singleCountry.Capital) > 0 {
		capital = singleCountry.Capital[0]
	}

	// Currency: pick the first currency key if available
	currency := ""
	if len(singleCountry.Currencies) > 0 {
		for key := range singleCountry.Currencies {
			currency = key
			break
		}
	}

	// Populate internal DTO
	country = dto.Country{
		Name:       singleCountry.Name.Common,
		Currency:   currency,
		Capital:    capital,
		Population: singleCountry.Population,
	}

	// ------------------------
	// Step 4: Cache the result for future requests
	// ------------------------
	cs.cacheService.Set(name, country)

	return country, nil
}

// GetCountryDataServiceCall calls the external REST Countries API and returns raw response.
func (cs CountryService) GetCountryDataServiceCall(ctx context.Context, name string) ([]dto.RestCountryAPIResponse, error) {
	var response []dto.RestCountryAPIResponse

	// Construct API endpoint
	endpoint := fmt.Sprintf(constants.COUNTRY_DETAIL_BY_NAME_API, cs.restCountryDomain, name)

	// Make GET request
	byteResponse, _, err := cs.sc.Get(ctx, endpoint, nil)

	if err != nil {
		log.Printf("Error sending GET request to '%s': %v\n", endpoint, err)
		return response, err
	}

	// Parse JSON response
	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		log.Printf("Error unmarshalling response from '%s': %v\n", endpoint, err)
		return response, err
	}

	return response, nil
}
