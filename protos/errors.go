package disposable

// Error - Get error message. Making it compatible with type error
func (e *Error) Error() string {
	return e.Message
}

// NewError - Will establish new error based on provided information.
// Function is here to provide unification layer.
func NewError(msg, errt string, err error) *Error {
	e := &Error{
		Message: msg,
		Type:    errt,
	}

	if err != nil {
		e.Info = map[string]string{
			"error": err.Error(),
		}
	}

	return e
}

// NewErrorFromInfo -
func NewErrorFromInfo(err *Error) *Error {
	e := &Error{
		Message: err.Message,
		Type:    err.Type,
	}

	if err.GetInfo() != nil {
		if errmsg, ok := err.GetInfo()["error"]; ok {
			e.Message = errmsg
		}
	}

	return e
}
