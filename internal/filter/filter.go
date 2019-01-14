package filter

import (
	"github.com/hcj112/wordfilter/internal/conf"
	"github.com/hcj112/wordfilter/internal/storage"
	"encoding/json"
	"github.com/hcj112/trie"
	"strings"
	"os"
	"io"
	"bufio"
)

type Filter struct {
	c    *conf.Config
	dao  *storage.Dao
	trie *trie.Trie
}

func New(c *conf.Config) *Filter {
	f := &Filter{
		c:    c,
		dao:  storage.New(c),
		trie: trie.New(),
	}
	go f.loadDict()
	go f.initBlackBucket()
	return f
}

func (f *Filter) Close() {
	f.dao.Close()
}

func (f *Filter) Add(keyword string) (err error) {
	err = f.dao.Storage.SaveKeyWord(json.RawMessage(keyword))
	if err != nil {
		return
	}

	if !f.trie.HasKeysWithPrefix(keyword) {
		f.trie.Add(keyword, json.RawMessage(keyword))
	}
	return
}

func (f *Filter) Remove(keyword string) (err error) {
	err = f.dao.Storage.DeleteKeyWord(json.RawMessage(keyword))
	if err != nil {
		return
	}
	f.trie.Remove(keyword)
	return
}

func (f *Filter) Filter(keyword string) (word string) {
	word = f.trie.Filter(keyword)
	return
}

func (f *Filter) List() (keywrods []string, err error) {
	keywrods, err = f.dao.Storage.BatchKeyWord()
	return
}

func (f *Filter) initBlackBucket() {
	keywrods, err := f.List()
	if err != nil {
		panic(err)
	}

	for _, keyword := range keywrods {
		err = f.Add(keyword)
		if err != nil {
			continue
		}
	}
}

func (f *Filter) loadDict() {
	if f.c.Dictionary.Path == "" {
		return
	}
	for _, file := range strings.Split(f.c.Dictionary.Path, ",") {
		dictFile, err := os.Open(file)
		defer dictFile.Close()

		if err != nil {
			panic(err)
		}

		buf := bufio.NewReader(dictFile)
		for {
			keyword, _, eof := buf.ReadLine()
			if eof == io.EOF {
				break
			}
			if string(keyword) == "" {
				continue
			}
			err := f.Add(string(keyword))
			if err != nil {
				continue
			}
		}
	}
}
