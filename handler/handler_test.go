package handler

import (
	"net/http"
	"net/http/httptest"
	"short_url/short"
	"short_url/short/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockShort struct {
	tinyFunc func(url model.Url) string
	getFunc func(tiny string) (string, error)
}

	// Tiny(url model.Url) string
	// Get(tinyUrl string) (string,error)
func (m *mockShort) Tiny(url model.Url) string {
	return m.tinyFunc(url)
}

func (m *mockShort) Get(tiny string) (string, error) {
	return m.getFunc(tiny)
}

func TestCreateTinyUrl(t *testing.T) {
	cases := []struct{
		name string
		request *http.Request
		response *httptest.ResponseRecorder
		httpStatusCodeExpected int
		msgStringExpected string
		tinyFunc func(url model.Url) string
		getFunc func(tiny string) (string, error)
	}{
		{
			name: "Test data post incomplete then return a error",
			request: httptest.NewRequest("POST", "/tiny", strings.NewReader(`{
				"long_url": "abcdefgh.com/long/dir"}`)),
			response: httptest.NewRecorder(),
			httpStatusCodeExpected: http.StatusBadRequest,
			msgStringExpected: "{\"msg\":\"user and long url are required\"}",
			tinyFunc: func(url model.Url) string {
				return ""
			},
			getFunc: func(tiny string) (string, error) {
				return "", nil
			},
		},
		{
			name: "Test len url is 0 then return a error",
			request: httptest.NewRequest("POST", "/tiny", strings.NewReader(`{"user_id": "user", "long_url": "a"}`)),
			response: httptest.NewRecorder(),
			httpStatusCodeExpected: http.StatusBadRequest,
			msgStringExpected: "{\"msg\":\"Url too short!\"}",
			tinyFunc: func(url model.Url) string {
				return ""
			},
			getFunc: func(tiny string) (string, error) {
				return "", nil
			},
		},
		{
			name: "Test len user is 0 then return a error",
			request: httptest.NewRequest("POST", "/tiny", strings.NewReader(`{"user_id": "a", "long_url": "abcdefgh.com/long"}`)),
			response: httptest.NewRecorder(),
			httpStatusCodeExpected: http.StatusBadRequest,
			msgStringExpected: "{\"msg\":\"User is empty\"}",
			tinyFunc: func(url model.Url) string {
				return ""
			},
			getFunc: func(tiny string) (string, error) {
				return "", nil
			},
		},
		{
			name: "Test data is ok then return tiny url",
			request: httptest.NewRequest("POST", "/tiny", strings.NewReader(`{"user_id": "yo", "long_url": "abcdefgh.com/long"}`)),
			response: httptest.NewRecorder(),
			httpStatusCodeExpected: http.StatusCreated,
			msgStringExpected: "{\"tiny_url\":\"abcdefg\"}",
			tinyFunc: func(url model.Url) string {
				return url.Long[0:7]
			},
			getFunc: func(tiny string) (string, error) {
				return "", nil
			},
		},
	}

	for _,c := range cases {
		t.Run(c.name, func(t *testing.T){
			c.request.Header.Add("Content-Type", "application/json")
			router := InitializeAndRun(nil, nil)	
			s = initMockShort(c.tinyFunc, c.getFunc)
			router.ServeHTTP(c.response, c.request)
			assert.Equal(t, c.httpStatusCodeExpected, c.response.Code)
			assert.Equal(t, c.msgStringExpected, c.response.Body.String())
		})
	}
}

func initMockShort(tinyFunc func(url model.Url) string, getFunc func(tiny string) (string, error)) short.Short{
	return &mockShort{
		tinyFunc: tinyFunc,
		getFunc: getFunc,
	}
}
