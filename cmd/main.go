package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/hcj112/wordfilter/internal/conf"
	"github.com/hcj112/wordfilter/internal/filter"
	"github.com/hcj112/wordfilter/internal/http"
)

const ver = "1.0.0"

func main() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		panic(err)
	}

	svr := filter.New(conf.Conf)
	httpSvr := http.New(conf.Conf.HTTPServer, svr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		glog.Infof("word filter get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			httpSvr.Close()
			glog.Warning("word filter [version: %s] exit", ver)
			glog.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

