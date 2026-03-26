package v1

import (
	"net/http"
	"search-country-system/internal/bootstrap"
)

func Router(app *bootstrap.App) *http.ServeMux {
	//Its responsible for the matching URL path with corresponding handler functions
	mux := http.NewServeMux()

	//Base path
	base := "/api/v1"

	mux.HandleFunc(base+"/countries/search", app.CountryHandler.GetCountryDetails)
	return mux
}
