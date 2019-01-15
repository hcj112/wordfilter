package conf

import (
	"flag"
	"time"
	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "example.toml", "default config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

type Config struct {
	Storage    *Storage
	Bolt       *Bolt
	Mysql      *Mysql
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
	DictPath    string
	KeywordPath string
}
