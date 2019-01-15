package filter

import (
	"github.com/hcj112/wordfilter/internal/conf"
	"github.com/hcj112/wordfilter/internal/storage"
	"encoding/json"
	"strings"
	"os"
	"io"
	"bufio"
	"github.com/huichen/sego"
	"unicode/utf8"
)

const REPLACE_STR = "*"

type Filter struct {
	c    *conf.Config
	dao  *storage.Dao
	sego sego.Segmenter
}

func New(c *conf.Config) *Filter {
	f := &Filter{
		c:   c,
		dao: storage.New(c),
	}
	go f.initDictionary()
	return f
}

func (f *Filter) Close() {
	f.dao.Close()
}

func (f *Filter) Add(keyword string) (err error) {
	key := f.getKey(keyword)
	err = f.dao.Storage.SaveKeyWord(key)
	return
}

func (f *Filter) Remove(keyword string) (err error) {
	key := f.getKey(keyword)
	err = f.dao.Storage.DeleteKeyWord(key)
	return
}

func (f *Filter) HasExist(keyword string) bool {
	key := f.getKey(keyword)
	if word, err := f.dao.Storage.GetKeyWord(key); word != "" && err == nil {
		return true
	}
	return false
}

func (f *Filter) getKey(keyword string) (key json.RawMessage) {
	key = json.RawMessage(strings.ToUpper(keyword))
	return
}

func (f *Filter) Filter(keyword string) string {
	bin := []byte(keyword)
	segments := f.sego.Segment(bin)
	keywords := make([]byte, 0, len(bin))
	for _, seg := range segments {
		word := bin[seg.Start():seg.End()]
		if f.HasExist(string(word)) {
			keywords = append(keywords, []byte(strings.Repeat(REPLACE_STR, utf8.RuneCount(word)))...)
		} else {
			keywords = append(keywords, word...)
		}
	}
	return string(keywords)
}

func (f *Filter) List() (keywrods []string, err error) {
	keywrods, err = f.dao.Storage.BatchKeyWord()
	return
}

func (f *Filter) initBlackBucket() {
	for _, file := range strings.Split(f.c.Dictionary.KeywordPath, ",") {
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
			word := string(keyword)
			if word == "" {
				continue
			}
			words := strings.Split(strings.ToUpper(strings.TrimSpace(word)), " ")
			if words[0] != "" {
				if err := f.Add(words[0]); err != nil {
					continue
				}
			}
		}
	}
}

/**
 * 导入分词库
 */
func (f *Filter) initDictionary() {
	f.sego.LoadDictionary(f.c.Dictionary.DictPath)
	f.initBlackBucket()
}
