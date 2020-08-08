package backup

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/assert"
)

func TestGetBackupSuccess(t *testing.T) {
	c := NewMysqlBackupController(&mockService{})

	next := openapi.NewRouter(c)
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

func TestGetBackupFail(t *testing.T) {
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
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestDeleteBackupSuccess(t *testing.T) {
	c := NewMysqlBackupController(&mockService{})

	next := openapi.NewRouter(c)
	r := httptest.NewRequest("DELETE", "/backup/bck", nil)

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

func TestDeleteBackupFail(t *testing.T) {
	c := &MysqlBackupController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}
	r := httptest.NewRequest("DELETE", "/backup/bck", nil)

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

func TestCreateBackupSuccess(t *testing.T) {
	c := NewMysqlBackupController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.Backup{
		Location:  "s3://bucket/loc/backup-1.dmp",
		Timestamp: time.Now(),
		S3access: openapi.S3Info{
			Bucket: "bucket",
			Path:   "/loc",
			AwsConfig: openapi.AwsConfig{
				AwsAccessKeyId:     "keyid",
				AwsSecretAccessKey: "secret",
				Region:             "us-east-1",
			},
		},
		Status:  "success",
		Message: "backup has succeeded",
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

func TestCreateBackupFailPayload(t *testing.T) {
	c := &MysqlBackupController{}

	next := openapi.NewRouter(c)
	c.service = &mockService{}

	data := []byte("payload error")

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
	assert.NotEqual(t, nil, err, "Should fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}

func TestCreateBackupFail(t *testing.T) {
	c := NewMysqlBackupController(&mockService{})

	next := openapi.NewRouter(c)

	b := openapi.Backup{
		Location:  "s3://bucket/loc/backup-1.dmp",
		Timestamp: time.Now(),
		S3access: openapi.S3Info{
			Bucket: "bucket",
			Path:   "/loc",
			AwsConfig: openapi.AwsConfig{
				AwsAccessKeyId:     "keyid",
				AwsSecretAccessKey: "secret",
				Region:             "us-east-1",
			},
		},
		Status:  "success",
		Message: "backup has succeeded",
	}
	data, err := json.Marshal(&b)

	r := httptest.NewRequest("POST", "/backup", bytes.NewReader(data))

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
	assert.NotEqual(t, nil, err, "Should Fail")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")
}
