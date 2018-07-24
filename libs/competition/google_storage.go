package competition

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gazoon/go-utils"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"path"
)

type GoogleFilesStorage struct {
	bucket     *storage.BucketHandle
	bucketName string
	httpClient *http.Client
}

func NewGoogleStorage(bucket string) (*GoogleFilesStorage, error) {
	ctx := context.Background()
	credsPath := path.Join(utils.GetCurrentFileDir(), "config", "google_cloud_keys.json")
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credsPath))
	if err != nil {
		return nil, errors.Wrap(err, "google storage client initialization")
	}
	httpClient := &http.Client{Timeout: httpTimeout}
	return &GoogleFilesStorage{
		bucket:     client.Bucket(bucket),
		bucketName: bucket,
		httpClient: httpClient}, nil
}

func (self *GoogleFilesStorage) Upload(ctx context.Context, storagePath, sourceUrl string) (string, error) {
	fileWriter := self.bucket.Object(storagePath).NewWriter(ctx)
	resp, err := self.httpClient.Get(sourceUrl)
	if err != nil {
		return "", errors.Wrap(err, "download source file")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("can't download the source, bad status code: %d", resp.StatusCode)
	}
	_, err = io.Copy(fileWriter, resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "copy content from source file to the storage")
	}
	err = fileWriter.Close()
	if err != nil {
		return "", errors.Wrap(err, "close storage writer")
	}
	return self.BuildUrl(storagePath), nil
}

func (self *GoogleFilesStorage) BuildUrl(path string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", self.bucketName, path)
}
