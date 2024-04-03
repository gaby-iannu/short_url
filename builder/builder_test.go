package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestBuild(t *testing.T) {
	repository, cache := Build()
	assert.NotNil(t, repository)
	assert.NotNil(t, cache)
}

func TestBuildeDataSource(t *testing.T) {
	datasource := buildDataSource()
	assert.NotNil(t, datasource)
	assert.NotEmpty(t, datasource.ConnMaxLifetime)
	assert.NotEmpty(t, datasource.DBName)
	assert.NotEmpty(t, datasource.Hostname)
	assert.NotEmpty(t, datasource.MaxIdleConns)
	assert.NotEmpty(t, datasource.MaxOpenConns)
	assert.NotEmpty(t, datasource.Password)
	assert.NotEmpty(t, datasource.Username)
}
