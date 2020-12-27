package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// NewStorage takes a S3 connection and creates a default storage
func NewStorage() *Storage {
	return &Storage{}
}

// Storage is the default storage for S3
type Storage struct {
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func getClient(ctx context.Context) (client *storage.Client, err error) {
	json := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if json != "" && isJSON(json) {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(json), storage.ScopeReadWrite)
		if err != nil {
			return nil, fmt.Errorf("google.JWTConfigFromJSON: %v", err)
		}
		ts := jwtConfig.TokenSource(ctx)
		return storage.NewClient(ctx, option.WithTokenSource(ts))
	}
	return storage.NewClient(ctx)
}

// Push pushes a file to blaqhole bucket
func (s *Storage) Push(request *openapi.BackupRequest, filename string) error {
	for _, v := range request.Envs {
		os.Setenv(v.Name, v.Value)
	}
	ctx := context.Background()
	client, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("Cannot get Client: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	location := request.Location
	if request.Location[0:1] == "/" {
		location = request.Location[1:]
	}
	wc := client.Bucket(request.Bucket).Object(location).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Printf("Blob %v uploaded.\n", location)
	return nil
}

// Pull pull a file from the blackhole
func (s *Storage) Pull(request *openapi.BackupRequest, filename string) error {
	for _, v := range request.Envs {
		os.Setenv(v.Name, v.Value)
	}
	ctx := context.Background()
	client, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("Cannot get Client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	location := request.Location
	if request.Location[0:1] == "/" {
		location = request.Location[1:]
	}
	rc, err := client.Bucket(request.Bucket).Object(location).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%s).NewReader: %v", location, err)
	}
	defer rc.Close()

	buf := make([]byte, 1024)
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fo.Close()

	for {
		// read a chunk
		n, err := rc.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	return nil
}

// Delete deletes a file from the blackhole
func (s *Storage) Delete(request *openapi.BackupRequest) error {
	for _, v := range request.Envs {
		os.Setenv(v.Name, v.Value)
	}
	ctx := context.Background()
	client, err := getClient(ctx)
	if err != nil {
		return fmt.Errorf("Cannot get Client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	location := request.Location
	if request.Location[0:1] == "/" {
		location = request.Location[1:]
	}
	o := client.Bucket(request.Bucket).Object(location)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", location, err)
	}
	fmt.Printf("Blob %v deleted.\n", location)
	return nil
}
