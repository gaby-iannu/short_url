package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestConfigureEnvironment_WhenScope_IsEmpty_ThenDefaultHost(t *testing.T) {
	shorturlHost, mysqlHost, redisHost := configureEnvironment()
	assert.Equal(t, "localhost", shorturlHost)
	assert.Equal(t, "127.0.0.1", mysqlHost)
	assert.Equal(t, "localhost", redisHost)
}

func TestConfigureEnvironment_WhenScope_IsProd_ThenProdHost(t *testing.T) {
	shortHost := "0.0.0.0"
	mysqlHost := "mysql-shorturl"
	redisHost := "redis-shorturl"
	scope := "prod"

	os.Setenv("SCOPE",scope)
	os.Setenv("SHORT_HOST", shortHost)
	os.Setenv("BD_HOST", mysqlHost)
	os.Setenv("CACHE_HOST", redisHost)

	s, m, r := configureEnvironment()
	assert.Equal(t, shortHost, s)
	assert.Equal(t, mysqlHost, m)
	assert.Equal(t, redisHost, r)
}

func TestConfigureEnvironment_WithErrors_ThenPanic(t *testing.T) {
	cases := []struct{
		name string
		scope string
		shortHost string
		mysqlHost string
		redisHost string
		err error
	}{
		{
			name: "Test Configure Environment When Scope Is Prod And shortHost Is Empty Then Panic",
			scope: "prod",
			shortHost: "",
			mysqlHost:"mysql-shorturl",
			redisHost: "redis-shorturl",
			err: fmt.Errorf("short host is empty"),
		},
		{
			name: "Test Configure Environment When Scope Is Prod And mysqlHost Is Empty Then Panic",
			scope: "prod",
			shortHost: "0.0.0.0",
			mysqlHost:"",
			redisHost: "redis-shorturl",
			err: fmt.Errorf("mysql host is empty"),
		},
		{
			name: "Test Configure Environment When Scope Is Prod And redisHost Is Empty Then Panic",
			scope: "prod",
			shortHost: "0.0.0.0",
			mysqlHost:"mysql-shorturl",
			redisHost: "",
			err: fmt.Errorf("redis host is empty"),
		},
	}

	for _,c := range cases {
		t.Run(c.name, func(t *testing.T){
			os.Setenv("SCOPE",c.scope)
			os.Setenv("SHORT_HOST", c.shortHost)
			os.Setenv("BD_HOST", c.mysqlHost)
			os.Setenv("CACHE_HOST", c.redisHost)
			assert.PanicsWithError(t, c.err.Error(), func(){configureEnvironment()})
		})
	}

}
