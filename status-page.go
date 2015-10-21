package jobrunner

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

type StatusData struct {
	Id        cron.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}

func StatusPage() []StatusData {

	ents := MainCron.Entries()

	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].Id = v.ID
		Statuses[k].JobRunner = AddJob(v.Job)
		Statuses[k].Next = v.Next
		Statuses[k].Prev = v.Prev

	}

	// t := template.New("status_page")

	// var data bytes.Buffer
	// t, _ = t.ParseFiles("views/Status.html")

	// t.ExecuteTemplate(&data, "status_page", Statuses())
	return Statuses
}

func StatusJson() map[string]interface{} {

	ents := MainCron.Entries()

	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].Id = v.ID
		Statuses[k].JobRunner = AddJob(v.Job)
		Statuses[k].Next = v.Next
		Statuses[k].Prev = v.Prev

	}

	return map[string]interface{}{
		"jobrunner": Statuses,
	}

}

func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
