package smolwordle

type corpusError string

func (e corpusError) Error() string {
	return string(e)
}

type gameError string

func (e gameError) Error() string {
	return string(e)
}
