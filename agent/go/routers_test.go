package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRouter struct{}

func (m *MockRouter) Routes() Routes {
	return Routes{
		{
			"CreateUser",
			strings.ToUpper("Post"),
			"/user",
			nil,
		},
	}
}

func TestRouter(t *testing.T) {
	m := &MockRouter{}
	r := NewRouter(m).GetRoute("CreateUser")
	methods, _ := r.GetMethods()
	assert.Equal(t, "POST", methods[0], "validate mux.Router")
}

func TestParseInt(t *testing.T) {
	v, err := parseIntParameter("64")
	assert.Equal(t, nil, err, "no errors")
	assert.Equal(t, int64(64), v, "validate int64")
}

func TestEncodeJSONResponse(t *testing.T) {
	w := httptest.NewRecorder()
	s := http.StatusInternalServerError
	EncodeJSONResponse(map[string]string{"name": "me"}, &s, w)
	response := w.Result()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u := map[string]string{}
	err = json.Unmarshal(bodyBytes, &u)
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "result http-500")

	w = httptest.NewRecorder()
	EncodeJSONResponse(map[string]string{"name": "me"}, nil, w)
	response = w.Result()
	bodyBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	u = map[string]string{}
	err = json.Unmarshal(bodyBytes, &u)
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, http.StatusOK, response.StatusCode, "result http-200")
}

func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

func TestReadFormFileToTempFile(t *testing.T) {
	b, w := createMultipartFormData(t, "image", "./model_s3_info.go")

	req, err := http.NewRequest("POST", "/", &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	_, err = ReadFormFileToTempFile(req, "xxx")
	assert.Error(t, err, "Should fail")
	_, err = ReadFormFileToTempFile(req, "image")
	assert.Error(t, err, "Should fail")

	os.MkdirAll("tmp", 0755)
	_, err = ReadFormFileToTempFile(req, "image")
	assert.Equal(t, nil, err, "Should succeed")
	os.RemoveAll("tmp/")
}

/*
// ReadFormFileToTempFile reads file data from a request form and writes it to a temporary file
func ReadFormFileToTempFile(r *http.Request, key string) (*os.File, error) {
        r.ParseForm()
        formFile, _, err := r.FormFile(key)
        if err != nil {
                return nil, err
        }

        defer formFile.Close()
        file, err := ioutil.TempFile("tmp", key)
        if err != nil {
                return nil, err
        }

        defer file.Close()
        fileBytes, err := ioutil.ReadAll(formFile)
        if err != nil {
                return nil, err
        }

        file.Write(fileBytes)
        return file, nil
}
*/
