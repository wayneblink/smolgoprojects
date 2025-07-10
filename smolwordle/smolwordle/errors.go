package smolwordle

type corpusError string

func (e corpusError) Error() string {
	return string(e)
}
