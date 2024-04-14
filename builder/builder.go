package builder

import (
	"fmt"
	"short_url/short/cache"
	"short_url/short/repository"
	"time"
)

func buildDataSource(mysqlHost string) repository.DataSource {
	return repository.DataSource{
		Username:"root",
		Password:"root",
		Hostname: fmt.Sprintf("%s:%s",mysqlHost, "3306"),
		DBName:"shorturl",
		MaxOpenConns:20,
		MaxIdleConns:20,
		ConnMaxLifetime: time.Minute * 3,
	}
}

func buildRepository(mysqlHost string) repository.Repository {
	return repository.New(buildDataSource(mysqlHost))
}

const (
	cache_port = 6379
)

func buildCache(redisHost string) cache.Cache {
	return cache.New(redisHost, cache_port)
}

// Build repository instance for DB persistences and
// build cache to get url faster
// mysqlHost: host to mysql can't be empty
// redisHost: host to redi can't be emptys
func Build(mysqlHost, redisHost string) (repository.Repository, cache.Cache, error) {
	if len(mysqlHost) == 0 {
		return nil, nil, fmt.Errorf("mysql host is empty")
	}

	if len(redisHost) == 0 {
		return nil, nil, fmt.Errorf("redis host is empty")
	}

	return buildRepository(mysqlHost), buildCache(redisHost), nil
}
