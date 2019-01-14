package storage

import (
	"github.com/hcj112/wordfilter/internal/conf"
	"github.com/golang/glog"
	"errors"
	"encoding/json"
)

var (
	useStorage     storage
	ErrStorageType = errors.New("unknown storage type")
)

const (
	BoltStorageType  = "bolt"
	//MySQLStorageType = "mysql"
)

type storage interface {
	SaveKeyWord(keyword json.RawMessage) error
	GetKeyWord(keyword json.RawMessage) (string, error)
	BatchKeyWord() ([]string, error)
	DeleteKeyWord(keyword json.RawMessage) error
	Close() error
}

func NewStorage(c *conf.Config) (storage, error) {
	if c.Storage.Type == BoltStorageType {
		useStorage = NewBoltStorage(c.Bolt)
	} else {
		glog.Errorf("unknown storage type: \"%s\"", c.Storage.Type)
		return nil, ErrStorageType
	}
	return useStorage, nil
}
