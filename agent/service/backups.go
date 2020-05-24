package service

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
	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

var (
	backups = make(map[time.Time]openapi.Backup)
	mutex   sync.Mutex
)

func InitializeBackup(b openapi.Backup) (*openapi.Backup, error) {
	if !b.Timestamp.IsZero() {
		return &b, errors.New("Timestamp already exists")
	}
	mutex.Lock()
	for {
		t := time.Now().Truncate(time.Second).UTC()
		if _, ok := backups[t]; !ok {
			b.Timestamp = t
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

func ExecuteBackup(b openapi.Backup) {
	t := b.Timestamp
	b.Status = "Running"
	mutex.Lock()
	backups[t] = b
	mutex.Unlock()
	filename := `backup-` + t.Format("20060102150405") + `.sql`
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
	err := PushS3File(filename, b.S3access.Credentials.AwsAccessKeyId, b.S3access.Credentials.AwsSecretAccessKey, b.S3access.Credentials.Region, b.S3access.Bucket, b.S3access.Path)
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

func PushS3File(filename, accesskey, secretkey, region, bucket, path string) error {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accesskey, secretkey, ""),
	})
	if err != nil {
		return err
	}
	file, err := os.Open("/tmp/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(path),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
