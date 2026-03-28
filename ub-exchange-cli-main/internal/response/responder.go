package response

import "net/http"

// APIResponse is the standard JSON envelope returned by all API endpoints.
type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success creates a successful APIResponse with HTTP 200 status.
func Success(data interface{}, message string) (resp APIResponse, status int) {
	resp = APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	return resp, http.StatusOK
}

// Error creates a failed APIResponse with the given message and HTTP status code.
func Error(message string, statusCode int, data interface{}) (resp APIResponse, status int) {
	if data == nil {
		data = make(map[string]string, 0)
	}
	resp = APIResponse{
		Status:  false,
		Message: message,
		Data:    data,
	}
	return resp, statusCode
}

// Unauthorized creates a failed APIResponse with HTTP 401 status.
func Unauthorized(data interface{}, message string) (resp APIResponse, status int) {
	resp = APIResponse{
		Status:  false,
		Message: message,
		Data:    data,
	}
	return resp, http.StatusUnauthorized
}
