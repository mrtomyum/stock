package api

type Response struct {
	Status  string
	Message string
	Result  interface{}
}

// Structure for collection of search string for frontend request.
type APISearch struct {
	Name string
}