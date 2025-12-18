package constants

// Code represents a business status code
type Code struct {
	Code    int
	Message string
}

var (
	Success       = Code{200, "success"}
	InvalidParams = Code{400, "invalid parameters"}
	InternalError = Code{500, "internal server error"}
)
