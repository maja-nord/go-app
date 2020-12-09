package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var conf = `
info: "/info"
`

type config struct {
	Info string
}

func main() {
	logrus.Info("Start application")

	c := &config{}
	if err := yaml.Unmarshal([]byte(conf), c); err != nil {
		logrus.Fatal(err)
	}

	r := chi.NewMux()
	r.Use(middleware.Timeout(time.Second * 10))

	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("serving /info")
		_, _ = w.Write([]byte(c.Info))
	})
	r.Get("/hash/{val}", func(w http.ResponseWriter, r *http.Request) {
		v := chi.URLParam(r, "val")
		sum, _ := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
		_,_=w.Write(sum[:])
	})

	_ = http.ListenAndServe(":8080", r)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}