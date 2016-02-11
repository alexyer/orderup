package main

// Custom bot errors.

// Error types
const (
	ARG_ERR = iota // Wrong command arguments error
	CMD_ERR        // Command logic error
)

type OrderupError struct {
	ErrType int
	msg     string
}

func NewOrderupError(msg string, errType int) *OrderupError {
	return &OrderupError{
		ErrType: errType,
		msg:     msg,
	}
}

func (e *OrderupError) Error() string {
	return e.msg
}
