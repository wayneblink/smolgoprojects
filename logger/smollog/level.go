package smollog

type Level byte

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)
