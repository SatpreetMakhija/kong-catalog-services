package datastore

type Service struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type ServiceSearchRequest struct {
	// Version field is used to filter based on exact version match. Eg: "v1.2.3"
	Version string   `json:"version"`
	Sort    []string `json:"sort"`
	// Name field is used to filter based on both exact term search as well as fuzzy search.
	Name    string   `json:"name"`
	// Query field represents natural language search. If Query field is non-empty,
	// Name, Version fields are ignored.
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type ServiceSearchResponse struct {
	Items    []Service `json:"items"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
