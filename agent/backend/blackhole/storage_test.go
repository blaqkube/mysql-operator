package blackhole

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gobuffalo/packr/v2"
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

	dumps := packr.New("dumps", "./dumps")
	blue, err := dumps.FindString("blue.sql")
	if err != nil {
		return err
	}
	c2 := string(dat)
	if c2 != blue {
		return errors.New("wrong file")
	}
	return nil
}

func deleteFile(filename string) error {
	return os.Remove(filename)
}

func (s *StorageSuite) TestBlackholeSuccess() {

	b := openapi.BackupRequest{
		Bucket:   os.Getenv("BACKUP_BUCKET"),
		Location: os.Getenv("BACKUP_LOCATION"),
		Envs:     []openapi.EnvVar{},
	}

	err := initFile("test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Push(&b, "test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = deleteFile("test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Pull(&b, "test2.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Delete(&b)
	assert.NoError(s.T(), err, "No Error")

	err = validFile("test2.txt", "test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = deleteFile("test2.txt")
	assert.NoError(s.T(), err, "No Error")

}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, &StorageSuite{})
}
