package storage

import (
	"github.com/hcj112/wordfilter/internal/conf"
)

type Dao struct {
	c       *conf.Config
	Storage storage
}

func New(c *conf.Config) *Dao {
	storage, err := NewStorage(c)
	if err != nil {
		panic(err)
	}
	d := &Dao{
		c:       c,
		Storage: storage,
	}
	return d
}

func (d *Dao) Close() {
	d.Storage.Close()
}


