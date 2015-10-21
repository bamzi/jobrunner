package jobrunner

import (
	"bytes"
	"text/template"

	"gopkg.in/robfig/cron.v2"
)

func StatusPage() *template.Template {
	t := template.New("status_page")
	t, _ = t.ParseFiles("views/Status.html")
	ent := MainCron.Entries()

	Entries := make([]cron.Entry, len(ent))

	for _, v := range ent {
		Entries = append(Entries, v)
	}
	var data bytes.Buffer
	t.ExecuteTemplate(&data, "status_page", Entries)
	return t
}
