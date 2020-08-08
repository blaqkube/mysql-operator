package user

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

func TestGetUserSuccess(t *testing.T) {
	c := NewMysqlUserController(&mockService{})
	next := openapi.NewRouter(c)
	r := httptest.NewRequest("GET", "/user/me", nil)

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
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result should succeed")
	assert.Equal(t, "me", u.Username, "Query Size should be me")
}

func TestGetUserFail(t *testing.T) {
	c := &MysqlUserController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/user/me", nil)

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

func TestGetUsers(t *testing.T) {
	c := &MysqlUserController{}
	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/user", nil)
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
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result http-200")
}

func TestGetUsersWithFailure(t *testing.T) {
	c := &MysqlUserController{}
	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/user", nil)
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
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestDeleteUserSuccess(t *testing.T) {
	c := NewMysqlUserController(&mockService{})

	next := openapi.NewRouter(c)
	r := httptest.NewRequest("DELETE", "/user/me", nil)

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

func TestDeleteUserFail(t *testing.T) {
	c := &MysqlUserController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("DELETE", "/user/me", nil)

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

func TestCreateUserSuccess(t *testing.T) {
	c := NewMysqlUserController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.User{
		Username: "me",
		Password: "me",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/user", bytes.NewReader(data))

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
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusCreated, response.StatusCode, "result should succeed")
	assert.Equal(t, "me", u.Username, "Query Size should be 1")
}

func TestCreateUserFailPayload(t *testing.T) {
	c := &MysqlUserController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}

	data := []byte("payload error")

	r := httptest.NewRequest("POST", "/user", bytes.NewReader(data))

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

func TestCreateUserFail(t *testing.T) {
	c := NewMysqlUserController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.User{
		Username: "me",
		Password: "me",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/user", bytes.NewReader(data))

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
