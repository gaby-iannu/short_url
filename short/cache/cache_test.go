package cache

import (
	"fmt"
	"log"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func setup() (cache, *miniredis.Miniredis) {
	mr, er := miniredis.Run()
	if er != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", er)
	}

	c := cache{
		client: redis.NewClient(&redis.Options{
				Addr: mr.Addr(),
			}),
	}

	return c, mr
}

func TestPut_ReturnError(t *testing.T) {
	c, mr := setup()

	err := fmt.Errorf("set error")
	mr.SetError(err.Error())
	defer mr.Close()


	e := c.Put("key","val")
	assert.EqualError(t, e, err.Error())
}

func TestPut_ErrorNil(t *testing.T) {
	c, mr := setup()
	defer mr.Close()

	mr.Set("key", "value")
	e := c.Put("key","value")
	assert.Nil(t, e)
}

func TestGet_ReturnError(t *testing.T) {
	c, mr := setup()
	defer mr.Close()

	mock := redismock.NewNiceMock(c.client)
	mock.On("Get", "key", "val", 0).
	Return(redis.NewStringCmd("", redis.Nil))

	_, err := c.Get("key")
	e := &NotExistError{}
	t.Log(err.Error())
	t.Log(e.Error())
	assert.EqualError(t, err, e.Error())
}

func TestGet_ReturnValue(t *testing.T) {
	c, mr := setup()
	defer mr.Close()

	mr.Set("key", "val")

	v, _ := c.Get("key")
	assert.Equal(t, "val", v)
}
