package jobrunner

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"os"
	"testing"
	"time"
)

type CustomJob struct {}

func (j CustomJob) Run() {
	fmt.Println("Custom job run")
}

type CustomLogger struct {
	cron.Logger
	Log *zerolog.Logger
}

func (l CustomLogger) Info(format string, keysAndValues ...interface{}) {
	l.Log.Info().Msgf(format, keysAndValues...)
}

func (l CustomLogger) Error(err error, format string, keysAndValues ...interface{}) {
	l.Log.Err(err).Msgf(format, keysAndValues...)
}

func Test_Init(t *testing.T) {
	zeroLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logger := CustomLogger{Log: &zeroLogger}

	StartWithLogger(logger)
	Schedule("@every 1s", CustomJob{})

	ch := make(chan bool)

	time.AfterFunc(5 * time.Second, func() {
		ch<-true
	})

	<-ch
}
