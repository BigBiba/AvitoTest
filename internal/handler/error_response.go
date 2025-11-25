package handler

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

var InternalError = ErrorResponse{
	Error: ErrorDetail{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
	},
}
