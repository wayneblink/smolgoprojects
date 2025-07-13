package repository

type Error string

func (r Error) Error() string {
	return string(r)
}

const ErrNotFound = Error("habit not found")
