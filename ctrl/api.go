package ctrl

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strings"
//)

type ResponseStatus int

const (
	SUCCESS ResponseStatus = iota
	FAIL
	ERROR
)

type Response struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	//Link    `json:"links"`
	Link
	Data interface{} `json:"data,omitempty"`
}

type Link struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
	First   string `json:"first,omitempty"`
	Prev    string `json:"prev,omitempty"`
	Next    string `json:"next,omitempty"`
	Last    string `json:"last,omitempty"`
}

// Structure for collection of search string for frontend request.
//type Search struct {
//	Name string
//}

//func (rs ResponseStatus) MarshalJSON() ([]byte, error) {
//	statusString, ok := map[ResponseStatus]string{
//		SUCCESS: "success",
//		FAIL:    "fail",
//		ERROR:   "error",
//	}[rs]
//	if !ok {
//		return nil, fmt.Errorf("invalid ResponseStatus value %v", rs)
//	}
//	return json.Marshal(statusString)
//}

//func (rs *ResponseStatus) UnmarshalJSON(data []byte) error {
//	// TODO: This method is not TEST yet!
//	// to receive response status from other service in JSON
//	// convert it to ENUM
//	var s string
//	if err := json.Unmarshal(data, &s); err != nil {
//		return fmt.Errorf("ResponseStatus should be a string, got %s", data)
//	}
//	s = strings.ToLower(s)
//	statusENUM, ok := map[string]ResponseStatus{
//		"success": SUCCESS,
//		"fail":    FAIL,
//		"error":   ERROR,
//	}[s]
//	if !ok {
//		return fmt.Errorf("invalid Response Status %q", s)
//	}
//	*rs = statusENUM
//	return nil
//}

//func (rs *Response) Success(data interface{}) {
//	var c gin.Context
//	rs.Status = SUCCESS
//	rs.Data = data
//	c.JSON(http.StatusOK, rs)
//}
