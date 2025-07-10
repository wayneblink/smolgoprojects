package ecbank

type ecbankError string

func (e ecbankError) Error() string {
	return string(e)
}

const (
	ErrCallingServer      = ecbankError("error calling server")
	ErrUnexpectedFormat   = ecbankError("unexpected response format")
	ErrChangeRateNotFound = ecbankError("couldn't find the exchange rate")
	ErrClientSide         = ecbankError("client side error when contacting ECB")
	ErrServerSide         = ecbankError("server side error when contacting ECB")
	ErrUnknownStatusCode  = ecbankError("error calling server")
)
