package jobrunner

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

const DEFAULT_JOB_POOL_SIZE = 10

var (
	// Singleton instance of the underlying job scheduler.
	MainCron *cron.Cron

	// This limits the number of jobs allowed to run concurrently.
	workPermits chan struct{}

	// Is a single job allowed to run concurrently with itself?
	selfConcurrent bool
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func Start(v ...int) {
	MainCron = cron.New()

	if len(v) > 0 {
		if v[0] > 0 {
			workPermits = make(chan struct{}, v[0])
		} else {
			workPermits = make(chan struct{}, DEFAULT_JOB_POOL_SIZE)
		}
	}

	if len(v) > 1 {
		if v[1] > 0 {
			selfConcurrent = true
		} else {
			selfConcurrent = false
		}
	}

	MainCron.Start()

	fmt.Printf("%s[JobRunner] %v Started... %s \n",
		magenta, time.Now().Format("2006/01/02 - 15:04:05"), reset)

}
