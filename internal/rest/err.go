package rest

type RestError struct {
	ErrorCode string                  `json:"errorCode"`
	Message   string                  `json:"message"`
	Data      *map[string]interface{} `json:"data,omitempty"`
}
