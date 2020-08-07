package backup

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/assert"
)

func TestGetBackupSuccess(t *testing.T) {
	c := &MysqlBackupController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/backup/bck", nil)

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
	assert.Equal(t, "s3://bucket/loc/backup-1.dmp", u.Location, "Query Size should be 1")
}

func TestGetBackupError(t *testing.T) {
	c := &MysqlBackupController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("GET", "/backup/bck", nil)

	r.Header.Set("apiKey", "test2")

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
	assert.Equal(t, http.StatusNotFound, response.StatusCode, "result should return Not Found")
}
