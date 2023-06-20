package error

type InternalError struct {
	Message string
}

func (internalError InternalError) Error() string {
	return internalError.Message
}
