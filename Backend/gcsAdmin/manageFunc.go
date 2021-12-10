package gcsAdmin
import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)


type ClientUploader struct {
	Client *storage.Client
	ProjectID string
	BucketName string
	UploadPath string
}

// delete file from gcs
func (c *ClientUploader) DeleteFile(filename string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	obj := c.Client.Bucket(c.BucketName).Object(c.UploadPath + filename)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("object(%q).delete: %v", filename, err)
	}
	return nil
}

// upload file to gcs
func (c *ClientUploader) UploadFile(file multipart.File, filename string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	writer := c.Client.Bucket(c.BucketName).Object(c.UploadPath + filename).NewWriter(ctx)
	writer.CacheControl = "no-store"
	writer.ContentType="image/png"

	if _, err := io.Copy(writer, file); err != nil {
		// gen err with format string
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("writer.Close: %v", err)
	}
	return nil
}