package conf

import (
	"flag"
	"time"
	"github.com/BurntSushi/toml"
)

var (
	confPath string
	dictPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "example.toml", "default config path")
	flag.StringVar(&dictPath, "dict", "dict.txt", "default dictionary path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if dictPath != "" {
		Conf.Dictionary.Path = dictPath
	}
	return
}

type Config struct {
	Storage    *Storage
	Bolt       *Bolt
	HTTPServer *HTTPServer
	Dictionary *Dictionary
}

type Storage struct {
	Type string
}

type Bolt struct {
	Path        string
	BlackBucket string
	DialTimeout time.Duration
}

type Mysql struct {
	Source string
}

type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Dictionary struct {
	Path string
}
