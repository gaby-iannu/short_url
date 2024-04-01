package short

import (
	"crypto/md5"
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



func toMd5(url model.Url) []byte {
	md5 := md5.New()
	io.WriteString(md5, url.Long)
	return md5.Sum(nil)
}

var varToMd5 func(model.Url) []byte = toMd5

// Reduce a long url and insert into DB
func (s *short) Tiny(url model.Url) string {
	return s.reduce(string(varToMd5(url)), url)
}



func (s *short) reduce(tiny string, url model.Url) string {
	
	if len(tiny) == 7 && s.repository.InsertIfNotExists(tiny, url) {
		return tiny
	}

	start := 0
	counter := 7
	l := len(tiny)
	for counter <= l && start < counter {
		urlReduced := tiny[start:counter]	
		if len(urlReduced) <= 7 && s.repository.InsertIfNotExists(urlReduced, url) {
			return urlReduced	
		}

		if (l - counter) >= 1 {
			counter++
		} 

		start++
	}

	return tiny
}
