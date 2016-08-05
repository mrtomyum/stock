package api

import (
	"encoding/json"
	"fmt"
)


type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
	Link    Link `json:"links,omitempty"`
}

type Link struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
	Next    string `json:"next,omitempty"`
	Last    string `json:"last,omitempty"`
}

// Structure for collection of search string for frontend request.
type Search struct {
	Name string
}

type ResponseStatus int

const (
	SUCCESS ResponseStatus = iota
	FAIL
	ERROR
)

func (rs ResponseStatus) MarshalJSON() ([]byte, error) {
	statusString, ok := map[ResponseStatus]string{
		SUCCESS: "success",
		FAIL:    "fail",
		ERROR:   "error",
	}[rs]
	if !ok {
		return nil, fmt.Errorf("invalid ResponseStatus value %v", rs)
	}
	return json.Marshal(statusString)
}
