package response

import "backend/config"

type Response struct {
	Code    config.ServiceCode `json:"code"`
	Data    map[string]any     `json:"data,omitempty"`
	Message string             `json:"message,omitempty"`
}

// NewResponse generates the response structure with the given service code
// and initializes the data map.
func New(code config.ServiceCode) *Response {
	return &Response{
		Code: code,
		Data: make(map[string]any),
	}
}

// AddKey adds the key to data map of the response with the given value
func (r *Response) AddKey(key string, value any) *Response {
	r.Data[key] = value
	return r
}

// SetMessage sets the value of message field in the response
func (r *Response) SetMessage(value string) *Response {
	r.Message = value
	return r
}

// NewErr generates the response with the given service code and error
func NewErr(code config.ServiceCode, err error) Response {
	return Response{
		Code:    code,
		Message: err.Error(),
	}
}
