package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cenkalti/backoff/v4"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
	maxRetries          = 3
)

// Object is a handle for a GCS object.
type Object struct {
	client   *s3.Client
	bucket   string
	name     string
	retry    bool
	metadata map[string]string
}

// NewObject creates a new Object.
func NewObject(client *s3.Client, bucket, name string) *Object {
	return &Object{
		client: client,
		bucket: bucket,
		name:   name,
		retry:  true,
	}
}

// Name returns the name of the object.
func (o *Object) Name() string {
	return o.name
}

// Bucket returns the name of the bucket.
func (o *Object) Bucket() string {
	return o.bucket
}

// SetRetry sets whether to retry operations on failure.
func (o *Object) SetRetry(retry bool) {
	o.retry = retry
}

// SetMetadata sets a custom metadata field on this object.
func (o *Object) SetMetadata(key, value string) {
	if o.metadata == nil {
		o.metadata = make(map[string]string)
	}

	o.metadata[key] = value
}

// Reader creates a reader for this object.
func (o *Object) Reader(ctx context.Context) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var reader io.ReadCloser
	operation := func() error {
		resp, err := o.client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(o.name),
		})
		if err != nil {
			return err
		}
		reader = resp.Body
		return nil
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	if err != nil {
		return nil, err
	}
	return reader, nil
}

// ReaderAt creates a reader that implements io.ReaderAt for this object.
// Note: This loads the entire object into memory.
func (o *Object) ReaderAt(ctx context.Context) (io.ReaderAt, error) {
	reader, err := o.Reader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read the entire object into memory
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to copy object to buffer: %w", err)
	}

	// Return a bytes.Reader which implements io.ReaderAt
	return bytes.NewReader(buf.Bytes()), nil
}

// Exists checks whether this object exists.
func (o *Object) Exists(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var exists bool
	operation := func() error {
		_, err := o.client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(o.name),
		})
		if err != nil {
			// The resource was deleted mid-flight, or never existed. Just treat as not found.
			return nil
		}
		exists = true
		return nil
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	if err != nil {
		return false, err
	}
	return exists, nil
}

// Delete removes this object.
func (o *Object) Delete(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	operation := func() error {
		_, err := o.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(o.name),
		})
		return err
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	return err
}

// Writer creates a writer for this object.
func (o *Object) Writer(ctx context.Context) (io.WriteCloser, error) {
	// Create a pipe to connect the writer to the upload
	pr, pw := io.Pipe()

	// Start a goroutine to upload the object using the pipe reader
	go func() {
		_, err := o.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(o.bucket),
			Key:         aws.String(o.name),
			Body:        pr,
			ContentType: aws.String("application/octet-stream"),
			Metadata:    o.metadata,
		})

		if err != nil {
			pr.CloseWithError(err)
		}
	}()

	return pw, nil
}

// ObjectAttr holds the attributes of an Object.
type ObjectAttr struct {
	// The name of the object.
	Name string

	// The size of the object in bytes.
	Size int64

	// The upload time of the object.
	Updated time.Time
}

// Attrs returns this object's attributes.
func (o *Object) Attrs(ctx context.Context) (*ObjectAttr, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var attr *ObjectAttr
	operation := func() error {
		resp, err := o.client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(o.name),
		})
		if err != nil {
			return err
		}

		attr = &ObjectAttr{
			Name: o.name,
		}

		if resp.ContentLength != nil {
			attr.Size = *resp.ContentLength
		}

		if resp.LastModified != nil {
			attr.Updated = *resp.LastModified
		}

		return nil
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	if err != nil {
		return nil, err
	}

	return attr, nil
}

// Size returns the size of this object.
func (o *Object) Size(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var size int64
	operation := func() error {
		resp, err := o.client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(o.name),
		})
		if err != nil {
			return err
		}
		if resp.ContentLength != nil {
			size = *resp.ContentLength
		}
		return nil
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	if err != nil {
		return 0, err
	}
	return size, nil
}

// Checksums returns the checksums for this object.
func (o *Object) Checksums(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	var checksum string
	operation := func() error {
		resp, err := o.client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(o.name),
		})
		if err != nil {
			return err
		}
		if resp.ETag != nil {
			checksum = *resp.ETag
		}
		return nil
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	if err != nil {
		return "", err
	}
	return checksum, nil
}

func convertTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// Create creates a new object.
func (o *Object) Create(ctx context.Context, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, defaultWriteTimeout)
	defer cancel()

	operation := func() error {
		_, err := o.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(o.bucket),
			Key:         aws.String(o.name),
			Body:        bytes.NewReader(data),
			ContentType: aws.String("application/octet-stream"),
			Metadata:    o.metadata,
		})
		return err
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	return err
}

// CopyFrom copies an object from another bucket and object.
func (o *Object) CopyFrom(ctx context.Context, srcBucket, srcObject string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultWriteTimeout)
	defer cancel()

	operation := func() error {
		_, err := o.client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(o.bucket),
			Key:        aws.String(o.name),
			CopySource: aws.String(fmt.Sprintf("%s/%s", srcBucket, srcObject)),
			Metadata:   o.metadata,
		})
		return err
	}

	var err error
	if o.retry {
		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)
		err = backoff.Retry(operation, b)
	} else {
		err = operation()
	}

	return err
}
