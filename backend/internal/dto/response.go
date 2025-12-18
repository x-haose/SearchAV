package dto

import "searchav/internal/model"

// Response is the unified response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
	List    interface{} `json:"list,omitempty"`
}

// WithCode sets the status code
func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

// WithMessage sets the message
func (r *Response) WithMessage(msg string) *Response {
	r.Message = msg
	return r
}

// WithData sets the data payload
func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

// WithList sets the list payload
func (r *Response) WithList(list interface{}) *Response {
	r.List = list
	return r
}

// SearchResponse is the search response structure
type SearchResponse struct {
	Code    int               `json:"code" example:"200"`
	Message string            `json:"msg" example:"success"`
	List    []model.VideoItem `json:"list"`
}

// DetailResponse is the detail response structure
type DetailResponse struct {
	Code    int               `json:"code" example:"200"`
	Message string            `json:"msg" example:"success"`
	Data    model.VideoDetail `json:"data"`
}

// ErrorResponse is the error response structure
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"msg" example:"bad request"`
}
