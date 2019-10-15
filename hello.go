package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
)

type config struct {
	Risky    bool
	Port     int
	Greeting string
}
type server struct {
	C *config
	P string
	V int
}

const (
	version  = 0
	tmplName = "homepage"
	tmplBody = `
{{define "` + tmplName + `" -}}
<html><body>
Version {{.V}} : {{if .C.Risky}}<em>{{end}}
{{- .C.Greeting}}{{if .C.Risky}}</em>{{end}} {{.P}}
</body></html>
{{end}}
`
)

var (
	enableRiskyFeature = flag.Bool("enableRiskyFeature", false,
		"Enables some risky feature.")
	port = flag.Int("port", 8080, "Port at which HTTP is served.")
)

func getConfig() *config {
	flag.Parse()
	greeting := os.Getenv("ALT_GREETING")
	if len(greeting) == 0 {
		greeting = "Hello"
	}
	return &config{*enableRiskyFeature, *port, greeting}
}

func main() {
	c := getConfig()
	tmpl := template.Must(template.New("main").Parse(tmplBody))
	http.HandleFunc("/quit",
		func(w http.ResponseWriter, r *http.Request) {
			go func() { time.Sleep(1 * time.Second); os.Exit(0) }()
		})
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if err := tmpl.ExecuteTemplate(
				w, tmplName, &server{c, r.URL.Path[1:], version}); err != nil {
				glog.Fatal(err)
			}
		})
	hostPort := ":" + strconv.Itoa(c.Port)
	if err := http.ListenAndServe(hostPort, nil); err != nil {
		glog.Fatal(err)
	}
}
