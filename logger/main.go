package main

import (
	"logger/smollog"
	"os"
	"time"
)

func main() {
	lgr := smollog.New(smollog.LevelInfo, smollog.WithOutput(os.Stdout))

	lgr.Infof("A little copying is better than a little dependency.")
	lgr.Errorf("Errors are values. Documentation is for %s.", "users")
	lgr.Debugf("Make the zero (%d) value useful.", 0)

	lgr.Infof("Hallo, %d %v", 2025, time.Now())
}
