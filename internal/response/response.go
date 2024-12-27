package response

type SuccessResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

type FailResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

type ErrorResponseOpts struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	ErrorResponseOpts
}

func NewSuccessResponse(status int, data any) SuccessResponse {
	return SuccessResponse{
		Status: status,
		Data:   data,
	}
}

func NewFailResponse(status int, data any) FailResponse {
	return FailResponse{
		Status: status,
		Data:   data,
	}
}

func NewErrorResponse(status int, message string, opts ErrorResponseOpts) ErrorResponse {
	return ErrorResponse{
		Status:            status,
		Message:           message,
		ErrorResponseOpts: opts,
	}
}
