package storage

import (
	"github.com/boltdb/bolt"
	"github.com/hcj112/wordfilter/internal/conf"
	"encoding/json"
	"time"
)

type BoltStorage struct {
	c  *conf.Bolt
	db *bolt.DB
}

func NewBoltStorage(c *conf.Bolt) *BoltStorage {
	options := bolt.Options{
		Timeout:  time.Duration(c.DialTimeout),
		ReadOnly: false,
	}
	db, err := bolt.Open(c.Path, 0600, &options)
	if err != nil {
		panic(err)
	}
	s := &BoltStorage{c: c, db: db}
	err = s.initBucket()
	if err != nil {
		panic(err)
	}
	return s
}

func (s *BoltStorage) SaveKeyWord(keyword json.RawMessage) (err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.c.BlackBucket))
		err = b.Put(keyword, keyword)
		return err
	})
	return
}

func (s *BoltStorage) GetKeyWord(key json.RawMessage) (keyword string, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.c.BlackBucket))
		keyword = string(b.Get(key))
		return nil
	})
	return
}

func (s *BoltStorage) BatchKeyWord() (keywords []string, err error) {
	keywords = []string{}
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.c.BlackBucket))
		err = b.ForEach(func(k, v []byte) error {
			keywords = append(keywords, string(v))
			return err
		})
		return err
	})
	return
}

func (s *BoltStorage) DeleteKeyWord(key json.RawMessage) (err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.c.BlackBucket))
		err := b.Delete(key)
		return err
	})
	return
}

func (s *BoltStorage) initBucket() (err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.c.BlackBucket))
		return err
	})
	return
}

func (s *BoltStorage) Close() error {
	return s.db.Close()
}
