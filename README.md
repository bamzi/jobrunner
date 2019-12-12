# ![](https://raw.githubusercontent.com/bamzi/jobrunner/master/views/runclock.jpg) JobRunner

JobRunner is framework for performing work asynchronously, outside of the request flow. It comes with cron to schedule and queue job functions for processing at specified time. 

It includes a live monitoring of current schedule and state of active jobs that can be outputed as JSON or Html template. 

## Install

`go get github.com/bamzi/jobrunner`

### Setup

```go
package main

import "github.com/bamzi/jobrunner"

func main() {
    jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
    jobrunner.Schedule("@every 5s", ReminderEmails{})
}

// Job Specific Functions
type ReminderEmails struct {
    // filtered
}

// ReminderEmails.Run() will get triggered automatically.
func (e ReminderEmails) Run() {
    // Queries the DB
    // Sends some email
    fmt.Printf("Every 5 sec send reminder emails \n")
}
```

### Live Monitoring
![](https://raw.githubusercontent.com/bamzi/jobrunner/master/views/jobrunner-html.png)
```go

// Example of GIN micro framework
func main() {
    routes := gin.Default()

    // Resource to return the JSON data
    routes.GET("/jobrunner/json", JobJson)

    // Load template file location relative to the current working directory
    routes.LoadHTMLGlob("../github.com/bamzi/jobrunner/views/Status.html")

    // Returns html page at given endpoint based on the loaded
    // template from above
    routes.GET("/jobrunner/html", JobHtml)

    routes.Run(":8080")
}

func JobJson(c *gin.Context) {
    // returns a map[string]interface{} that can be marshalled as JSON
    c.JSON(200, jobrunner.StatusJson())
}

func JobHtml(c *gin.Context) {
    // Returns the template data pre-parsed
    c.HTML(200, "", jobrunner.StatusPage())

}

```
## Custom Logger
If you use your own or customized logger package, you can integrate your logger package by following the steps below.

```go
package main

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

func logger1() {
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

func logger2() {
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

func logger3() {
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

func main() {
    logger1()
    logger2()
    logger3()
}
```
Log Output
```
// default Logger
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 JobRunner Started
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 start
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 added%!(EXTRA string=now, time.Time=2019-12-12 14:42:00.871014911 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-12 14:42:01 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 wake%!(EXTRA string=now, time.Time=2019-12-12 14:42:01.000175995 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 run%!(EXTRA string=now, time.Time=2019-12-12 14:42:01.000175995 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-12 14:42:02 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 wake%!(EXTRA string=now, time.Time=2019-12-12 14:42:02.000453377 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 run%!(EXTRA string=now, time.Time=2019-12-12 14:42:02.000453377 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-12 14:42:03 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 wake%!(EXTRA string=now, time.Time=2019-12-12 14:42:03.000297283 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 run%!(EXTRA string=now, time.Time=2019-12-12 14:42:03.000297283 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-12 14:42:04 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 wake%!(EXTRA string=now, time.Time=2019-12-12 14:42:04.000363516 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 run%!(EXTRA string=now, time.Time=2019-12-12 14:42:04.000363516 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-12 14:42:05 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 wake%!(EXTRA string=now, time.Time=2019-12-12 14:42:05.00030873 +0300 +03)
2019-12-12 11:42:00.870926568 +0000 UTC2019/12/12 run%!(EXTRA string=now, time.Time=2019-12-12 14:42:05.00030873 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-12 14:42:06 +0300 +03)

// zerolog.Logger
{"level":"info","time":"2019-12-11T19:19:11+03:00","message":"JobRunner Started"}
{"level":"info","time":"2019-12-11T19:19:11+03:00","message":"start"}
{"level":"info","time":"2019-12-11T19:19:11+03:00","message":"added%!(EXTRA string=now, time.Time=2019-12-11 19:19:11.539557059 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-11 19:19:12 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:12+03:00","message":"wake%!(EXTRA string=now, time.Time=2019-12-11 19:19:12.000175121 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:12+03:00","message":"run%!(EXTRA string=now, time.Time=2019-12-11 19:19:12.000175121 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-11 19:19:13 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:13+03:00","message":"wake%!(EXTRA string=now, time.Time=2019-12-11 19:19:13.000324711 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:13+03:00","message":"run%!(EXTRA string=now, time.Time=2019-12-11 19:19:13.000324711 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-11 19:19:14 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:14+03:00","message":"wake%!(EXTRA string=now, time.Time=2019-12-11 19:19:14.000597951 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:14+03:00","message":"run%!(EXTRA string=now, time.Time=2019-12-11 19:19:14.000597951 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-11 19:19:15 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:15+03:00","message":"wake%!(EXTRA string=now, time.Time=2019-12-11 19:19:15.0004445 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:15+03:00","message":"run%!(EXTRA string=now, time.Time=2019-12-11 19:19:15.0004445 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-11 19:19:16 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:16+03:00","message":"wake%!(EXTRA string=now, time.Time=2019-12-11 19:19:16.000527496 +0300 +03)"}
{"level":"info","time":"2019-12-11T19:19:16+03:00","message":"run%!(EXTRA string=now, time.Time=2019-12-11 19:19:16.000527496 +0300 +03, string=entry, cron.EntryID=1, string=next, time.Time=2019-12-11 19:19:17 +0300 +03)"}

// zap.Logger
2019-12-12T14:42:36.135+0300	INFO	jobrunner/init_test.go:101	JobRunner Started
2019-12-12T14:42:36.135+0300	INFO	jobrunner/init_test.go:101	start
2019-12-12T14:42:36.136+0300	INFO	jobrunner/init_test.go:101	added	{"now": "2019-12-12T14:42:36.136+0300", "entry": 1, "next": "2019-12-12T14:42:37.000+0300"}
2019-12-12T14:42:37.000+0300	INFO	jobrunner/init_test.go:101	wake	{"now": "2019-12-12T14:42:37.000+0300"}
2019-12-12T14:42:37.000+0300	INFO	jobrunner/init_test.go:101	run	{"now": "2019-12-12T14:42:37.000+0300", "entry": 1, "next": "2019-12-12T14:42:38.000+0300"}
2019-12-12T14:42:38.000+0300	INFO	jobrunner/init_test.go:101	wake	{"now": "2019-12-12T14:42:38.000+0300"}
2019-12-12T14:42:38.000+0300	INFO	jobrunner/init_test.go:101	run	{"now": "2019-12-12T14:42:38.000+0300", "entry": 1, "next": "2019-12-12T14:42:39.000+0300"}
2019-12-12T14:42:39.000+0300	INFO	jobrunner/init_test.go:101	wake	{"now": "2019-12-12T14:42:39.000+0300"}
2019-12-12T14:42:39.000+0300	INFO	jobrunner/init_test.go:101	run	{"now": "2019-12-12T14:42:39.000+0300", "entry": 1, "next": "2019-12-12T14:42:40.000+0300"}
2019-12-12T14:42:40.000+0300	INFO	jobrunner/init_test.go:101	wake	{"now": "2019-12-12T14:42:40.000+0300"}
2019-12-12T14:42:40.000+0300	INFO	jobrunner/init_test.go:101	run	{"now": "2019-12-12T14:42:40.000+0300", "entry": 1, "next": "2019-12-12T14:42:41.000+0300"}
2019-12-12T14:42:41.000+0300	INFO	jobrunner/init_test.go:101	wake	{"now": "2019-12-12T14:42:41.000+0300"}
2019-12-12T14:42:41.000+0300	INFO	jobrunner/init_test.go:101	run	{"now": "2019-12-12T14:42:41.000+0300", "entry": 1, "next": "2019-12-12T14:42:42.000+0300"}

```

## WHY?
To reduce our http response latency by 200+%

JobRunner was created to help us process functions unrelated to response without any delays to the http response. GoRoutines would timeout because response has already been processed, closing the instance window all together. 

Instead of creating a separate independent app, we wanted to save time and manageability of our current app by coupling-in the job processing. We did not want to have micro services. It's premature optimization.

If you have a web app or api service backend and want a job processing framework built into your app then JobRunner is for you. Once you hit mass growth and need to scale, you can simply decouple you JobRunners into a dedicated app.

## Use cases
Here are some examples of what we use JobRunner for:

* Send emails to new users after signup
* Sending push notification or emails based on specifics
* ReMarketing Engine - send invites, reminder emails, etc ...
* Clean DB, data or AMZ S3
* Sending Server stats to monitoring apps
* Send data stats at daily or weekly intervals

### Supported Featured
*All jobs are processed outside of the request flow*

* Now: process a job immediately
* In: processing a job one time, after a given time
* Every: process a recurring job after every given time gap
* Schedule: process a job (recurring or otherwise) at a given time


## Compatibility

JobRunner is designed to be framework agnostic. So it will work with pure Go apps as well as any framework written in Go Language. 

*Verified Supported Frameworks*

* Gin
* Echo
* Martini
* Beego
* Revel (JobRunner is modified version of revel's Jobs module)
* ...

**Examples & recipes are coming soon**

## Basics

```go
    jobrunner.Schedule("* */5 * * * *", DoSomething{}) // every 5min do something
    jobrunner.Schedule("@every 1h30m10s", ReminderEmails{})
    jobrunner.Schedule("@midnight", DataStats{}) // every midnight do this..
    jobrunner.Every(16*time.Minute, CleanS3{}) // evey 16 min clean...
    jobrunner.In(10*time.Second, WelcomeEmail{}) // one time job. starts after 10sec
    jobrunner.Now(NowDo{}) // do the job as soon as it's triggered
```
[**More Detailed CRON Specs**](https://github.com/robfig/cron/blob/v2/doc.go)

## Contribute

**Use issues for everything**

- Report problems
- Discuss before sending pull request
- Suggest new features
- Improve/fix documentation

## Credits
- [revel jobs module](https://github.com/revel/modules/tree/master/jobs) - Origin of JobRunner
- [robfig cron v3](https://github.com/robfig/cron/tree/v3) - github.com/robfig/cron/v3
- [contributors](https://github.com/bamzi/jobrunner/graphs/contributors)

### Author 
**Bam Azizi**
Github: [@bamzi](https://github.com/bamzi)
Twitter: [@bamazizi](https://twitter/bamazizi)

#### License
MIT