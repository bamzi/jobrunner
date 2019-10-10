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

	functions =[]interface{}{makeWorkPermits,isSelfConcurrent}
)

func makeWorkPermits(bufferCapacity int) {
	if bufferCapacity <=0 {
		workPermits = make(chan struct{}, DEFAULT_JOB_POOL_SIZE)
	} else {
		workPermits = make(chan struct{}, bufferCapacity)
	}
}

func isSelfConcurrent(cocnurrencyFlag int) {
	if cocnurrencyFlag <=0 {
		selfConcurrent = false
	} else {
		selfConcurrent = true
	}
}

func Start(v ...int) {
	MainCron = cron.New()

	for i,option := range v {
		functions[i].(func(int))(option)
	}


	MainCron.Start()

	fmt.Printf("%s[JobRunner] %v Started... %s \n",
		magenta, time.Now().Format("2006/01/02 - 15:04:05"), reset)

}
