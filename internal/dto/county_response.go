package dto

type Country struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int    `json:"population"`
}

type RestCountryAPIResponse struct {
	Name       Name     `json:"name"`
	Capital    []string `json:"capital"`
	Population int      `json:"population"`
	Currencies map[string]struct {
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

type Name struct {
	Common string `json:"common"`
}
