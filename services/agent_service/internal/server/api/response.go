package api

// ============================================================================
// Generic API Response
// ============================================================================

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func SuccessResponse[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Data:    data,
	}
}

func SuccessMessageResponse[T any](message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse[any] {
	return APIResponse[any]{
		Success: false,
		Message: message,
	}
}
