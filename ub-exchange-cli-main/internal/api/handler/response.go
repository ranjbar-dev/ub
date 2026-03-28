package handler

// APIResponse is the standard envelope for all API responses.
// All endpoints must return this format for consistency.
type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse creates an error response with optional validation field errors.
func NewErrorResponse(message string, data map[string]string) APIResponse {
	return APIResponse{
		Status:  false,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponse creates a success response with optional data payload.
func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
}
