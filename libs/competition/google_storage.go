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

type googleStorage struct {
	bucket     *storage.BucketHandle
	bucketName string
	httpClient *http.Client
}

func newGoogleStorage(bucket string) (*googleStorage, error) {
	ctx := context.Background()
	credsPath := path.Join(utils.GetCurrentFileDir(), "config", "google_cloud_keys.json")
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credsPath))
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Timeout: httpTimeout}
	return &googleStorage{
		bucket:     client.Bucket(bucket),
		bucketName: bucket,
		httpClient: httpClient}, nil
}

func (self *googleStorage) upload(ctx context.Context, storagePath, sourceUrl string) (string, error) {
	fileWriter := self.bucket.Object(storagePath).NewWriter(ctx)
	resp, err := http.Get(sourceUrl)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("can't download the source, bad status code: %d", resp.StatusCode)
	}
	_, err = io.Copy(fileWriter, resp.Body)
	if err != nil {
		return "", err
	}
	err = fileWriter.Close()
	if err != nil {
		return "", err
	}
	return self.buildUrl(storagePath), nil
}

func (self *googleStorage) buildUrl(path string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", self.bucketName, path)
}
