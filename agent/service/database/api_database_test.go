package database

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

func TestGetDatabaseSuccess(t *testing.T) {
	c := NewMysqlDatabaseController(&mockService{})
	next := openapi.NewRouter(c)
	r := httptest.NewRequest("GET", "/database/me", nil)

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Database{}
	err = json.Unmarshal(bodyBytes, u)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result should succeed")
	assert.Equal(t, "me", u.Name, "Query Size should be me")
}

func TestGetDatabaseFail(t *testing.T) {
	c := &MysqlDatabaseController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/database/me", nil)

	r.Header.Set("apiKey", "test2")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Database{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestDeleteDatabaseSuccess(t *testing.T) {
	c := NewMysqlDatabaseController(&mockService{})

	next := openapi.NewRouter(c)
	r := httptest.NewRequest("DELETE", "/database/me", nil)

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Message{}
	err = json.Unmarshal(bodyBytes, u)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result should succeed")
}

func TestDeleteDatabaseFail(t *testing.T) {
	c := &MysqlDatabaseController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("DELETE", "/database/me", nil)

	r.Header.Set("apiKey", "test2")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Message{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestCreateDatabaseSuccess(t *testing.T) {
	c := NewMysqlDatabaseController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.Database{
		Name: "me",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/database", bytes.NewReader(data))

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Database{}
	err = json.Unmarshal(bodyBytes, u)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusCreated, response.StatusCode, "result should succeed")
	assert.Equal(t, "me", u.Name, "Query Size should be 1")
}

func TestCreateDatabaseFailPayload(t *testing.T) {
	c := &MysqlDatabaseController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}

	data := []byte("payload error")

	r := httptest.NewRequest("POST", "/database", bytes.NewReader(data))

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Database{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestCreateDatabaseFail(t *testing.T) {
	c := NewMysqlDatabaseController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.Database{
		Name: "me",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/database", bytes.NewReader(data))

	r.Header.Set("apiKey", "test2")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Database{}
	err = json.Unmarshal(bodyBytes, u)
	assert.NotEqual(t, nil, err, "Should Fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}
