package controllers

import (
	"strings"

	"github.com/bamzi/jobrunner"
	"github.com/revel/revel"
	"github.com/robfig/cron"
)

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Status() revel.Result {
	remoteAddress := c.Request.RemoteAddr
	if revel.Config.BoolDefault("jobs.acceptproxyaddress", false) {
		if proxiedAddress, isProxied := c.Request.Header["X-Forwarded-For"]; isProxied {
			remoteAddress = proxiedAddress[0]
		}
	}
	if !strings.HasPrefix(remoteAddress, "127.0.0.1") && !strings.HasPrefix(remoteAddress, "::1") {
		return c.Forbidden("%s is not local", remoteAddress)
	}
	entries := jobrunner.MainCron.Entries()
	return c.Render(entries)
}

func init() {
	revel.TemplateFuncs["castjob"] = func(job cron.Job) *jobrunner.Job {
		return job.(*jobrunner.Job)
	}
}
