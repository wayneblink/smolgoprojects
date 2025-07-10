package session

type Error string

func (d Error) Error() string {
	return string(d)
}

const (
	ErrGameOver = Error("❌game over")
)
