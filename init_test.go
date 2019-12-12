package jobrunner

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

type CustomJob struct {}

func (j CustomJob) Run() {
	fmt.Println("Custom job run with default logger")
}

type CustomLogger struct {
	cron.Logger
	Log *log.Logger
}

func (l CustomLogger) Info(format string, keysAndValues ...interface{}) {
	l.Log.Printf(format, keysAndValues...)
}

func (l CustomLogger) Error(err error, format string, keysAndValues ...interface{}) {
	l.Log.Println(err)
	l.Log.Printf(format, keysAndValues...)
}

func Test_InitOne(t *testing.T) {
	now := time.Now().UTC()
	defaultLogger := log.New(os.Stderr, now.String(), 1)

	logger := CustomLogger{Log: defaultLogger}

	StartWithLogger(logger)
	Schedule("@every 1s", CustomJob{})

	ch := make(chan bool)

	time.AfterFunc(5 * time.Second, func() {
		ch<-true
	})

	<-ch
}

type CustomJob2 struct {}

func (j CustomJob2) Run() {
	fmt.Println("Custom job run with zerolog.Logger")
}

type CustomLogger2 struct {
	cron.Logger
	Log *zerolog.Logger
}

func (l CustomLogger2) Info(format string, keysAndValues ...interface{}) {
	l.Log.Info().Msgf(format, keysAndValues...)
}

func (l CustomLogger2) Error(err error, format string, keysAndValues ...interface{}) {
	l.Log.Err(err).Msgf(format, keysAndValues...)
}

func Test_InitTwo(t *testing.T) {
	zeroLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logger := CustomLogger2{Log: &zeroLogger}

	StartWithLogger(logger)
	Schedule("@every 1s", CustomJob2{})

	ch := make(chan bool)

	time.AfterFunc(5 * time.Second, func() {
		ch<-true
	})

	<-ch
}

type CustomJob3 struct {}

func (j CustomJob3) Run() {
	fmt.Println("Custom job run with zap.Logger")
}

type CustomLogger3 struct {
	cron.Logger
	Log *zap.Logger
}

func (l CustomLogger3) Info(format string, keysAndValues ...interface{}) {
	l.Log.Sugar().Infow(format, keysAndValues...)
}

func (l CustomLogger3) Error(err error, format string, keysAndValues ...interface{}) {
	l.Log.Sugar().Errorw(strings.Join([]string{err.Error(), format}, ""), keysAndValues...)
}

func Test_InitThree(t *testing.T) {
	zapLogger, _ := zap.NewDevelopment()

	logger := CustomLogger3{Log: zapLogger}

	StartWithLogger(logger)
	Schedule("@every 1s", CustomJob3{})

	ch := make(chan bool)

	time.AfterFunc(5 * time.Second, func() {
		ch<-true
	})

	<-ch
}