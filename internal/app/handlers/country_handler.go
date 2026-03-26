package handlers

import (
	"net/http"
	"search-country-system/internal/app/services/interfaces"
)

// CountryHandler handles HTTP requests related to country data.
// It embeds BaseHandler to reuse context and JSON response methods.
type CountryHandler struct {
	BaseHandler
	cs interfaces.ICountryService // Service interface to fetch country data
}

// InitCountryHandler initializes a new CountryHandler with a given country service.
func InitCountryHandler(cs interfaces.ICountryService) *CountryHandler {
	return &CountryHandler{
		cs: cs, // Assign the service to the handler
	}
}

// GetCountryDetails handles HTTP GET requests to fetch details of a country by name.
// Example URL: /country?name=India
func (ch *CountryHandler) GetCountryDetails(w http.ResponseWriter, r *http.Request) {
	// Extract the "name" query parameter from the URL
	name := r.URL.Query().Get("name")
	if name == "" {
		// If no name is provided, respond with HTTP 400 Bad Request
		ch.ResponseError(w, http.StatusBadRequest, "name query param is required")
		return
	}

	// Use the request's context
	ctx := r.Context()

	// Call the country service to fetch country data
	result, err := ch.cs.GetCountryData(ctx, name)
	if err != nil {
		// If the service returns an error, respond with HTTP 500 Internal Server Error
		ch.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// On success, respond with HTTP 200 OK and the country data
	ch.ResponseOK(w, result)
}
