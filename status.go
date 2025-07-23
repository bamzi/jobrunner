package jobrunner

import (
	"time"

	"github.com/robfig/cron/v3"
)

type StatusData struct {
	Id        cron.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}

// Return detailed list of currently running recurring jobs
// to remove an entry, first retrieve the ID of entry
func Entries() []cron.Entry {
	return MainCron.Entries()
}

func StatusPage() []StatusData {
	ents := MainCron.Entries()

	Statuses := make([]StatusData, 0, len(ents))
	for _, v := range ents {
		Statuses = append(Statuses, StatusData{
			Id:        v.ID,
			JobRunner: AddJob(v.Job),
			Next:      v.Next,
			Prev:      v.Prev,
		})
	}

	// t := template.New("status_page")

	// var data bytes.Buffer
	// t, _ = t.ParseFiles("views/Status.html")

	// t.ExecuteTemplate(&data, "status_page", Statuses())
	return Statuses
}

func StatusJson() map[string]interface{} {
	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}
}

func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
