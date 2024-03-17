package short

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"short_url/short/cache"
	"short_url/short/model"
	"short_url/short/repository"
)

type Short interface {
	Tiny(url model.Url) string
	Get(tinyUrl string) (string,error)
}

type short struct {
	repository repository.Repository
	cache cache.Cache
}

// type ShortOption func(*Short) 

func New(cache cache.Cache, repository repository.Repository) Short {
	return &short{
		repository: repository,
		cache: cache,
	}
}

// Get long url
func (s *short) Get(tinyUrl string) (string,error) {
	
	longUrl,err := s.cache.Get(tinyUrl)
	if  _, ok := err.(*cache.NotExistError); ok {
		url:=s.repository.Read(tinyUrl)
		if url == (model.Url{}) {
			return "", fmt.Errorf("url dosen't exit")
		}
		s.cache.Put(tinyUrl, url.Long)
		return url.Long, nil
	}

	return longUrl, nil
}

// Reduce a long url and insert into DB
func (s *short) Tiny(url model.Url) string {
	return s.toMd5(url)
}

func (s *short) toMd5(url model.Url) string {
	md5 := md5.New()
	io.WriteString(md5, url.Long)
	b := md5.Sum(nil)
	return s.reduce(hex.EncodeToString(b), url)
}

func (s *short) reduce(tiny string, url model.Url) string {
	
	if len(tiny) == 7 && s.repository.InsertIfNotExists(tiny, url) {
		return tiny
	}

	start := 0
	counter := 7
	l := len(tiny)
	for counter < l {
		urlReduced := tiny[start:counter]	
		if len(urlReduced) <= 7 && s.repository.InsertIfNotExists(urlReduced, url) {
			return urlReduced	
		}
		start++
		counter++
	}


	return ""
}
