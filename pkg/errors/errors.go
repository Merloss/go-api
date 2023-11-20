package errors

type ErrorCode = string

const (
	INTERNAL_SERVER_ERROR ErrorCode = "INTERNAL_SERVER_ERROR"
	NOT_FOUND             ErrorCode = "NOT_FOUND"
	BAD_REQUEST           ErrorCode = "BAD_REQUEST"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	RequestId string    `json:"requestId,omitempty"`
}
