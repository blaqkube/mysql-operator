package grant

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/assert"
)

func TestGetGrantSuccess(t *testing.T) {
	c := NewMysqlGrantController(&mockService{})
	next := openapi.NewRouter(c)
	r := httptest.NewRequest("GET", "/user/me/database/pong/grant", nil)

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Grant{}
	err = json.Unmarshal(bodyBytes, u)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result should succeed")
	assert.Equal(t, "readWrite", u.AccessMode, "Query Size should be me")
}

func TestGetUserFail(t *testing.T) {
	c := &MysqlGrantController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/user/me/database/pong/grant", nil)

	r.Header.Set("apiKey", "test2")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.User{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestCreateGrantSuccess(t *testing.T) {
	c := NewMysqlGrantController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.Grant{
		AccessMode: "readWrite",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/user/me/database/pong/grant", bytes.NewReader(data))

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	g := &openapi.Grant{}
	err = json.Unmarshal(bodyBytes, g)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusCreated, response.StatusCode, "result should succeed")
	assert.Equal(t, ReadWriteAccessMode, g.AccessMode, "Query Size should be 1")
}

func TestCreateGrantFailPayload(t *testing.T) {
	c := &MysqlGrantController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}

	data := []byte("payload error")

	r := httptest.NewRequest("POST", "/user/me/database/pong/grant", bytes.NewReader(data))

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.User{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestCreateGrantFail(t *testing.T) {
	c := NewMysqlGrantController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.Grant{
		AccessMode: "readWrite",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/user/me/database/pong/grant", bytes.NewReader(data))

	r.Header.Set("apiKey", "test2")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.User{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should Fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}
