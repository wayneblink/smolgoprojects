package money

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidDecimal = Error("unable to convert the decimal")
	ErrTooLarge       = Error("quantity over 10^12 is too large")
)
