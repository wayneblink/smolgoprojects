package smolwordle

import "strings"

type hint int

const (
	absentCharacter hint = iota
	wrongPosition
	correctPosition
)

func (h hint) String() string {
	switch h {
	case absentCharacter:
		return "-"
	case wrongPosition:
		return "?"
	case correctPosition:
		return "+"
	default:
		return "💔"
	}
}

type Feedback []hint

func (fb Feedback) String() string {
	sb := strings.Builder{}
	for _, s := range fb {
		sb.WriteString(s.String())
	}
	return sb.String()
}

func (fb Feedback) Equal(other Feedback) bool {
	if len(fb) != len(other) {
		return false
	}
	for index, value := range fb {
		if value != other[index] {
			return false
		}
	}
	return true
}

func (fb Feedback) GameWon() bool {
	for _, c := range fb {
		if c != correctPosition {
			return false
		}
	}

	return true
}
