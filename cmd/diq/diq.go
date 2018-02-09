package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(os.Stderr)

	logrus.Info("Hello world!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("got a request!")
		fmt.Fprintf(w, "Hello world!")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Errorf("Problem listening: %s", err)
	}
}
