package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/datainq/go-workshop/pkg/topic"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 1024
)

func main() {
	var tmplDir, staticDir, addr string
	flag.StringVar(&tmplDir, "tmplDir", "./tmpl", "template directory")
	flag.StringVar(&staticDir, "staticDir", "./static",
		"directory with static resource (JS, CSS)")
	flag.StringVar(&addr, "addr", ":8080", "address to listen at")
	flag.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(os.Stderr)

	logrus.Info("Hello world!")

	upgrader := &websocket.Upgrader{
		ReadBufferSize:    socketBufferSize,
		WriteBufferSize:   socketBufferSize,
		EnableCompression: true,
	}
	tp := topic.NewTopic("", "")
	go tp.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("got a ws request!")
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.Errorf("failed to upgrade connection: %s", err)
		}
		client := topic.NewClient(socket, tp)
		go client.ReadPipe()
		client.WritePipe()
	})
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
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Errorf("Problem listening: %s", err)
	}
}
