package main

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

var (
	tmplDir   = "./tmpl"
	staticDir = "./static"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(os.Stderr)

	logrus.Info("Hello world!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("got a request!")
		f, err := os.Open(path.Join(tmplDir, "index.html"))
		if err != nil {
			logrus.Errorf("cannot open file: %s", err)
		}
		defer f.Close()
		if _, err = io.Copy(w, f); err != nil {
			logrus.Errorf("problem copying HTML: %s", err)
		}
	})
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Errorf("Problem listening: %s", err)
	}
}
