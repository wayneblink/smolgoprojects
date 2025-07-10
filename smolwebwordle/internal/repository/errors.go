package repository

type Error string

func (r Error) Error() string {
	return string(r)
}

const (
	ErrNotFound      = Error("game not found in repository")
	ErrConflictingID = Error("cannot create game with already-existing ID")
)
