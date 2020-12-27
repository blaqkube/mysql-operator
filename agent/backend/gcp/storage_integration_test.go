// +build integration

package gcp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type StorageSuite struct {
	suite.Suite
	Storage *Storage
}

func (s *StorageSuite) SetupTest() {
	s.Storage = NewStorage()
}

func initFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s]", filename)
	return nil
}

func validFile(filename, content string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	c1 := fmt.Sprintf("[test.txt]")
	c2 := string(dat)
	if c2 != c1 {
		return errors.New("wrong file")
	}
	return nil
}

func deleteFile(filename string) error {
	return os.Remove(filename)
}

func (s *StorageSuite) TestGCPStorageSuccess() {

	_ = godotenv.Load()

	b := openapi.BackupRequest{
		Bucket:   os.Getenv("BACKUP_BUCKET"),
		Location: os.Getenv("BACKUP_LOCATION"),
		Envs: []openapi.EnvVar{
			{
				Name:  "GOOGLE_APPLICATION_CREDENTIALS",
				Value: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
			},
		},
	}

	err := initFile("test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Push(&b, "test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Pull(&b, "test2.txt")
	assert.NoError(s.T(), err, "No Error")

	err = deleteFile("test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Delete(&b)
	assert.NoError(s.T(), err, "No Error")

	err = validFile("test2.txt", "test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = deleteFile("test2.txt")
	assert.NoError(s.T(), err, "No Error")
}

func (s *StorageSuite) TestGCPStorageFail() {

	_ = godotenv.Load()

	b := openapi.BackupRequest{
		Bucket:   "failed",
		Location: os.Getenv("BACKUP_LOCATION"),
		Envs: []openapi.EnvVar{
			{
				Name:  "GOOGLE_APPLICATION_CREDENTIALS",
				Value: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
			},
		},
	}

	c := openapi.BackupRequest{
		Bucket:   os.Getenv("BACKUP_BUCKET"),
		Location: os.Getenv("BACKUP_LOCATION"),
		Envs: []openapi.EnvVar{
			{
				Name:  "GOOGLE_APPLICATION_CREDENTIALS",
				Value: "unknown",
			},
		},
	}

	err := initFile("test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Push(&b, "test.txt")
	assert.Error(s.T(), err, "Error")
	assert.Regexp(s.T(), "Writer.Close.*", err.Error(), "Writer.Close")

	err = s.Storage.Push(&c, "test.txt")
	assert.Error(s.T(), err, "Error")
	assert.Regexp(s.T(), "Cannot get Client.*", err.Error(), "Cannot get Client")

	err = s.Storage.Push(&b, "test2.txt")
	assert.Error(s.T(), err, "Error")
	assert.Regexp(s.T(), "open.*", err.Error(), "AccessDenied")

	err = s.Storage.Pull(&b, "test2.txt")
	assert.Error(s.T(), err, "Error")
	assert.Regexp(s.T(), "Object.*", err.Error(), "AccessDenied")

	err = deleteFile("test.txt")
	assert.NoError(s.T(), err, "No Error")

}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, &StorageSuite{})
}
