package builder

import (
	"short_url/short/cache"
	"short_url/short/repository"
	"time"
)

func buildDataSource() repository.DataSource {
	return repository.DataSource{
		Username:"root",
		Password:"root",
		Hostname:"127.0.0.1:3306",
		DBName:"shorturl",
		MaxOpenConns:20,
		MaxIdleConns:20,
		ConnMaxLifetime: time.Minute * 3,
	}
}

func buildRepository() repository.Repository {
	return repository.New(buildDataSource())
}

const (
	cache_url = "localhost"
	cache_port = 6379
)

func buildCache() cache.Cache {
	return cache.New(cache_url, cache_port)
}

// Build repository instance for DB persistences and
// build cache to get url faster
func Build() (repository.Repository, cache.Cache) {
	return buildRepository(), buildCache()
}
