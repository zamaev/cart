package customerror

type ErrStatusCode struct {
	msg    string
	Status int
}

func (e ErrStatusCode) Error() string {
	return e.msg
}

func NewErrStatusCode(msg string, status int) ErrStatusCode {
	return ErrStatusCode{msg, status}
}
