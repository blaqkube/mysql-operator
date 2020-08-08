package mysql

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

var (
	backups = make(map[time.Time]openapi.Backup)
	mutex   sync.Mutex
)

// S3MysqlBackup defines an interface for the backup tools
type S3MysqlBackup interface {
	InitializeBackup(openapi.Backup) (*openapi.Backup, error)
	ExecuteBackup(b openapi.Backup)
	GetBackup(t time.Time) (*openapi.Backup, error)
}

// NewS3MysqlBackup takes a S3 connection and creates a default backup
func NewS3MysqlBackup() S3MysqlBackup {
	return &S3MysqlDefaultBackup{}
}

// S3MysqlDefaultBackup is the default S3MysqlBackup
type S3MysqlDefaultBackup struct {
}

// GetBackup returns a backup from the execution time
func (s *S3MysqlDefaultBackup) GetBackup(t time.Time) (*openapi.Backup, error) {
	b, ok := backups[t]
	if !ok {
		return nil, errors.New("no backup")
	}
	return &b, nil
}

// InitializeBackup registers a backup to start it later
func (s *S3MysqlDefaultBackup) InitializeBackup(b openapi.Backup) (*openapi.Backup, error) {
	if !b.Timestamp.IsZero() {
		return &b, errors.New("Timestamp already exists")
	}
	mutex.Lock()
	for {
		t := time.Now().Truncate(time.Second).UTC()
		if _, ok := backups[t]; !ok {
			b.Timestamp = t
			filename := `backup-` + t.Format("20060102150405") + `.sql`
			b.Location = "s3://" + b.S3access.Bucket + b.S3access.Path + "/" + filename
			b.Status = "Pending"
			backups[t] = b
			mutex.Unlock()
			return &b, nil
		}
		mutex.Unlock()
		time.Sleep(1 * time.Second)
		mutex.Lock()
	}
}

// ExecuteBackup executes an initialized backup
func (s *S3MysqlDefaultBackup) ExecuteBackup(b openapi.Backup) {
	t := b.Timestamp
	filename := `backup-` + t.Format("20060102150405") + `.sql`
	b.Status = "Running"
	mutex.Lock()
	backups[t] = b
	mutex.Unlock()
	cmd := exec.Command("mysqldump", "--all-databases", "--lock-all-tables", "--host=127.0.0.1", `--result-file=/tmp/`+filename)
	if err := cmd.Run(); err != nil {
		b.Status = "Failed"
		b.Message = fmt.Sprintf("%v", err)
		mutex.Lock()
		backups[t] = b
		mutex.Unlock()
		return
	}
	b.Status = "Pushing to S3"
	mutex.Lock()
	backups[t] = b
	mutex.Unlock()
	err := s.PushS3File(&b, filename)
	b.Status = "Available"
	if err != nil {
		b.Status = "Failed"
		b.Message = fmt.Sprintf("%v", err)
	}
	mutex.Lock()
	backups[t] = b
	mutex.Unlock()
	return
}

// PushS3File pushes a file to S3
func (s *S3MysqlDefaultBackup) PushS3File(backup *openapi.Backup, filename string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(backup.S3access.AwsConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			backup.S3access.AwsConfig.AwsAccessKeyId,
			backup.S3access.AwsConfig.AwsSecretAccessKey,
			"",
		),
	})
	if err != nil {
		return err
	}
	file, err := os.Open("/tmp/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(backup.S3access.Bucket),
		Key:                aws.String(backup.S3access.Path + "/" + filename),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}

// PullS3File pull a file from S3, using a different location if necessary
func (s *S3MysqlDefaultBackup) PullS3File(backup *openapi.Backup, location, filename string) error {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	l := location
	if location == "" {
		l = backup.S3access.Path
	}
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(backup.S3access.Bucket),
			Key:    aws.String(l),
		})
	return err
}
