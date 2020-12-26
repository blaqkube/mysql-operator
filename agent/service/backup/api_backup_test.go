package backup

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

func TestGetBackupByIDSuccess(t *testing.T) {
	c := NewController(&mockService{})

	next := openapi.NewRouter(c)
	r := httptest.NewRequest("GET", "/backup/abcd", nil)

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Backup{}
	err = json.Unmarshal(bodyBytes, u)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result should succeed")
	assert.Equal(t, "/loc/backup-1.dmp", u.Location, "Should return location")
}

func TestGetBackupByIDFailure(t *testing.T) {
	c := NewController(&mockService{})

	next := openapi.NewRouter(c)
	r := httptest.NewRequest("GET", "/backup/abcd", nil)

	r.Header.Set("apiKey", "test2")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	_, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result should succeed")
}

func TestCreateBackupSuccess(t *testing.T) {
	c := NewController(&mockService{})

	next := openapi.NewRouter(c)
	b := openapi.BackupRequest{
		Backend:  "s3",
		Location: "/bucket/loc/backup-1.dmp",
		Bucket:   "bucket",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/backup", bytes.NewReader(data))

	r.Header.Set("apiKey", "test1")

	w := httptest.NewRecorder()
	next.ServeHTTP(w, r)
	response := w.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := &openapi.Backup{}
	err = json.Unmarshal(bodyBytes, u)
	assert.Equal(t, err, nil, "Should succeed")
	assert.Equal(t, http.StatusCreated, response.StatusCode, "result should succeed")
	assert.Equal(t, "s3://bucket/loc/backup-1.dmp", u.Location, "Query Size should be 1")
}
