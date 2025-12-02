package assets_api

type responseAPI struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Status  interface{} `json:"status"`
	Data    interface{} `json:"result,omitempty"`
	Paging  interface{} `json:"paging,omitempty"`
	Filter  interface{} `json:"filters,omitempty"`
}

// Response is a general structure for API responses an array data.
func SuccessResponseData(message string, data, paging, filter interface{}) *responseAPI {
	return &responseAPI{Status: "success", Data: data, Code: 200, Message: message, Paging: paging, Filter: filter}
}

// Response is a general structure for API responses a data.
func SimpSuccessResponse(message string, data interface{}) *responseAPI {
	return SuccessResponseData(message, data, nil, nil)
}
func BadRequestResponse(message string) *responseAPI {
	return &responseAPI{Status: "error", Code: 400, Message: message}
}

func UnauthenticationResponse() *responseAPI {
	return &responseAPI{Status: "authentication", Code: 401, Message: "Unauthenticated request"}
}
func UnAuthorizationResponse() *responseAPI {
	return &responseAPI{Status: "forbidden", Code: 403, Message: "Forbidden request"}
}
func NotFoundResponse(message string) *responseAPI {
	return &responseAPI{Status: "notfound", Code: 404, Message: message}
}
func ResponseError(code int, message ...string) *responseAPI {
	defaultMessages := map[int]string{
		400: "Unknown error",
		401: "Unauthenticated request",
		403: "Forbidden request",
		404: "Not found request",
	}

	// Định nghĩa status mặc định dựa trên mã lỗi
	statusMap := map[int]string{
		400: "error",
		401: "authentication",
		403: "forbidden",
		404: "notfound",
	}
	finalMessage := defaultMessages[code]
	if len(message) > 0 && message[0] != "" {
		finalMessage = message[0]
	}
	return &responseAPI{
		Code:    code,
		Status:  statusMap[code],
		Message: finalMessage,
	}
}
