//go:generate go-enum -f=$GOFILE

package logging

import (
"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
)

var once sync.Once

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableLevelTruncation: true, FullTimestamp:true})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.WithError(err).Warning("could not parse log level, so using default.")
	} else {
		logrus.SetLevel(lvl)
	}
}

// Trace prints the calling function file, name and line.
func Trace(traceType TraceType) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	logrus.Tracef("%s %s,:%d %s\n", traceType.String(), frame.File, frame.Line, frame.Function)
}

// TraceType explains if the trace is entering a function, is exciting a function or if its still in a function.
/*
ENUM(
entering
inside
exiting
)
*/
type TraceType int
