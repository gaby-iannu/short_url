package short

import (
	"fmt"
	"short_url/short/model"
	"testing"

	"short_url/short/cache"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	insertIfNotExistsFunc func(tiny string, url model.Url) bool
	readFunc func(tiny string) model.Url 
}

func (r *mockRepository) InsertIfNotExists(tiny string, url model.Url) bool {
	return r.insertIfNotExistsFunc(tiny, url)
}

func (r *mockRepository) Read(tinyUrl string) model.Url {
	return r.readFunc(tinyUrl)
}

type mockCache struct {
	putFunc func(key, value string) error
	getFunc func(key string) (string, error)
}

func (c *mockCache) Put(key, value string) error {
	return c.putFunc(key, value)
}

func (c *mockCache) Get(key string) (string, error) {
	return c.getFunc(key)
}

func TestReduce(t *testing.T) {
	cases := []struct{
		name string
		insertIfNotExistsFunc func(tiny string, url model.Url) bool
		readFunc func(tiny string) model.Url
		putFunc func(key, value string) error
		getFunc func(key string) (string, error)
		url model.Url
		tiny string
		expected string

	}{
		{
			name: "Test reduce len of tiny equal 7 and InsertIfNotExists true",
			insertIfNotExistsFunc: func(tiny string, url model.Url) bool {
				return true
			},
			readFunc: func(tiny string) model.Url{
				return model.Url{}
			},
			putFunc: func(key string, value string) error {
				return nil
			},
			getFunc: func(key string) (string, error) {
				return "", nil
			},
			url: model.Url{},
			tiny: "1234567",
			expected: "1234567",
		},
		{
			name: "Test reduce len of tiny equal 7 and InsertIfNotExists false then return tiny",
			insertIfNotExistsFunc: func(tiny string, url model.Url) bool {
				return false
			},
			readFunc: func(tiny string) model.Url{
				return model.Url{}
			},
			putFunc: func(key string, value string) error {
				return nil
			},
			getFunc: func(key string) (string, error) {
				return "", nil
			},
			url: model.Url{},
			tiny: "1234567",
			expected: "1234567",
		},
		{
			name: "Test reduce len of tiny equal 6 and InsertIfNotExists false then return tiny",
			insertIfNotExistsFunc: func(tiny string, url model.Url) bool {
				return false
			},
			readFunc: func(tiny string) model.Url{
				return model.Url{}
			},
			putFunc: func(key string, value string) error {
				return nil
			},
			getFunc: func(key string) (string, error) {
				return "", nil
			},
			url: model.Url{},
			tiny: "123456",
			expected: "123456",
		},
		{
			name: "Test reduce len of tiny equal 8 and InsertIfNotExists true then return tiny reduced",
			insertIfNotExistsFunc: func(tiny string, url model.Url) bool {
				return true
			},
			readFunc: func(tiny string) model.Url{
				return model.Url{}
			},
			putFunc: func(key string, value string) error {
				return nil
			},
			getFunc: func(key string) (string, error) {
				return "", nil
			},
			url: model.Url{},
			tiny: "12345678",
			expected: "1234567",
		},
		{
			name: "Test iterate tiny url until get an URL with 7 characters long",
			insertIfNotExistsFunc: func(tiny string, url model.Url) bool {
				if tiny == "5678910" {
					return true
				}

				return false
			},
			readFunc: func(tiny string) model.Url{
				return model.Url{}
			},
			putFunc: func(key string, value string) error {
				return nil
			},
			getFunc: func(key string) (string,error) {
				return "", nil
			},
			url: model.Url{},
			tiny: "12345678910",
			expected: "5678910",
		},
	}

	for _,c := range cases {
		t.Run(c.name, func(t *testing.T){
			s := buildShort(c.insertIfNotExistsFunc, c.readFunc, c.putFunc, c.getFunc)
			assert.Equal(t, c.expected, s.reduce(c.tiny, c.url))
		})
	}
}

/*
*/
func TestTiny(t *testing.T) {

	cases := []struct{
		name string
		insertIfNotExistsFunc func(tiny string, url model.Url) bool
		readFunc func(tiny string) model.Url
		putFunc func(key, value string) error
		getFunc func(key string) (string, error)
		url model.Url
		expected string

	}{
		{
			name: "Test facebook url",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(string) (string, error) {
				return "", nil
			},
			url: model.Url{
				Long: "https://www.facebook.com/",
				User: "user",
			},
			expected: "e203e98",
		},
		{
			name: "Test instagram url",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(string) (string, error) {
				return "", nil
			},
			url: model.Url{
				Long: "https://www.instagram.com/",
				User: "user",
			},
			expected: "bcbafb6",
		},
		{
			name: "Test gmail url",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(string) (string, error) {
				return "", nil
			},
			url: model.Url{
				Long: "https://mail.google.com/mail/u/0/#inbox", 
				User: "user",
			},
			expected: "2122c56",
		},
		{
			name: "Test go play url",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(string) (string, error) {
				return "", nil
			},
			url: model.Url{
				Long: "https://go.dev/play/",
				User: "user",
			},
			expected: "f111d2b",
		},

	}

	for _,c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := buildShort(c.insertIfNotExistsFunc, c.readFunc, c.putFunc, c.getFunc)
			encodeToString = func([]byte) string {
				return c.expected
			}
			assert.Equal(t, c.expected, s.Tiny(c.url))
		})
	}
}

func TestGet(t *testing.T) {

	cases := []struct{
		name string
		insertIfNotExistsFunc func(tiny string, url model.Url) bool
		readFunc func(tiny string) model.Url
		putFunc func(key, value string) error
		getFunc func(key string) (string,error)
		tiny string
		expectedValue string
		expectedError error

	}{
		{
			name: "Test cache return error distinct to NotExistsError then return empty string and nil",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(string)(string, error) {
				return "", fmt.Errorf("other error")
			},
			tiny: "abc",
			expectedValue: "",
			expectedError: nil,
		},
		{
			name: "Test cache return value then return value and nil",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(key string)(string, error) {
				if key == "abc" {
					return "abcdefg.com/dir", nil
				} 
				return "", nil
			},
			tiny: "abc",
			expectedValue: "abcdefg.com/dir",
			expectedError: nil,
		},
		{
			name: "Test cache return NotExistsError, repository return empty Url  then return empty string and error",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(key string)(string, error) {
				return "", &cache.NotExistError{}
			},
			tiny: "abc",
			expectedValue: "",
			expectedError: fmt.Errorf("url dosen't exit"),
		},
		{
			name: "Test cache return NotExistsError, repository return Url then return long url and nil",
			insertIfNotExistsFunc: func(string, model.Url) bool {
				return true
			},
			readFunc: func(string) model.Url {
				return model.Url{
					Long: "abcdefg.com/dir",
					User: "user",
				}
			},
			putFunc: func(string, string) error {
				return nil
			},
			getFunc: func(key string)(string, error) {
				return "", &cache.NotExistError{}
			},
			tiny: "abc",
			expectedValue: "abcdefg.com/dir",
			expectedError: nil,
		},
	
	}

	for _,c := range cases {
		t.Run(c.name, func(t *testing.T){
			s := buildShort(c.insertIfNotExistsFunc, c.readFunc, c.putFunc, c.getFunc)

			longUrl, err := s.Get(c.tiny)
			assert.Equal(t, c.expectedValue, longUrl)
			if c.expectedError == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, c.expectedError.Error())
			}
		})
	}
}

func buildShort(insertIfNotExistsFunc func(string, model.Url) bool, readFunc func(string) model.Url, 
						putFunc func(string, string) error, getFunc func(string) (string,error))  *short {

	ss := New(&mockCache{putFunc: putFunc, getFunc: getFunc}, 
		&mockRepository{insertIfNotExistsFunc: insertIfNotExistsFunc, readFunc: readFunc})
	 return ss.(*short)
}
