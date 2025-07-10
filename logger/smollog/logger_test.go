package smollog_test

import (
	"logger/smollog"
	"testing"
)

type testWriter struct {
	contents string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.contents = tw.contents + string(p)
	return len(p), nil
}

const (
	debugMessage = "Why write I still all one, ever the same,"
	infoMessage  = "And keep invention in a noted weed,"
	errorMessage = "That every word doth almost tell my name,"
)

func TestLogger_DebufInfoErrorf(t *testing.T) {
	type testCase struct {
		level    smollog.Level
		expected string
	}

	tt := map[string]testCase{
		"debug": {
			level:    smollog.LevelDebug,
			expected: debugMessage + "\n" + infoMessage + "\n" + errorMessage + "\n",
		},
		"info": {
			level:    smollog.LevelInfo,
			expected: infoMessage + "\n" + errorMessage + "\n",
		},
		"error": {
			level:    smollog.LevelError,
			expected: errorMessage + "\n",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tw := &testWriter{}

			testedLogger := smollog.New(tc.level, smollog.WithOutput(tw))

			testedLogger.Debugf(debugMessage)
			testedLogger.Infof(infoMessage)
			testedLogger.Errorf(errorMessage)

			if tw.contents != tc.expected {
				t.Errorf("invalid contents, expected %q, got %q", tc.expected, tw.contents)
			}
		})
	}
}
