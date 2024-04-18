package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestBuild(t *testing.T) {
	repository, cache, err  := Build("mysqlHost","redisHost")
	assert.NotNil(t, repository)
	assert.NotNil(t, cache)
	assert.Nil(t, err)
}

func TestBuild_WhenMysqlHost_IsNull_ThenReturnError(t *testing.T) {
	repository, cache, err := Build("","")
	assert.Nil(t, repository)
	assert.Nil(t, cache)
	assert.EqualError(t, err, "mysql host is empty")
}

func TestBuild_WhenRedisHost_IsNull_ThenReturnError(t *testing.T) {
	repository, cache, err := Build("mysqlHost","")
	assert.Nil(t, repository)
	assert.Nil(t, cache)
	assert.EqualError(t, err, "redis host is empty")
}

func TestBuildeDataSource(t *testing.T) {
	datasource := buildDataSource("mysqlHost")
	assert.NotNil(t, datasource)
	assert.NotEmpty(t, datasource.ConnMaxLifetime)
	assert.NotEmpty(t, datasource.DBName)
	assert.NotEmpty(t, datasource.Hostname)
	assert.NotEmpty(t, datasource.MaxIdleConns)
	assert.NotEmpty(t, datasource.MaxOpenConns)
	assert.NotEmpty(t, datasource.Password)
	assert.NotEmpty(t, datasource.Username)
}
