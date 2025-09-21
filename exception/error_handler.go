package exception

type ErrorHandler struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorHandler) Error() string {
	return e.Message
}

func NotFound(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    404,
		Message: message,
	}
}
func BadRequest(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    400,
		Message: message,
	}
}
func InternalServerError(message string) *ErrorHandler {
	return &ErrorHandler{
		Code:    500,
		Message: message,
	}
}
