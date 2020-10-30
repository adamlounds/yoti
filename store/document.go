package store

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type myS3 interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type DocumentS3Store struct {
	s3Client myS3
}
type DocumentStore interface {
	Retrieve(id []byte) (payload []byte, err error)
	Store(id, payload []byte) (err error)
}

func NewDocumentS3Store() (*DocumentS3Store, error) {
	s3Config := &aws.Config{
		Credentials:  credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		Endpoint: aws.String("http://minio:9000"),
		Region: aws.String("eu-west-1"),
		DisableSSL: aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	s, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(s)

	return &DocumentS3Store{s3Client}, nil
}

func (d *DocumentS3Store) Store(id, payload []byte) (err error) {
	resp, err := d.s3Client.PutObject(&s3.PutObjectInput{
		Body: bytes.NewReader(payload),
		Bucket: aws.String("storage"),
		Key: aws.String(hex.EncodeToString(id)),
	})
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "uploaded %v", resp.ETag)

	return err
}

func (d *DocumentS3Store) Retrieve(id []byte) (payload []byte, err error) {
	//buf := &aws.WriteAtBuffer{}
	//downloader := s3manager.NewDownloader(d.s3Client.(client.ConfigProvider))
	//numBytes, err := downloader.Download(buf, &s3.GetObjectInput{
	//	Bucket: aws.String("storage"),
	//	Key: aws.String(hex.EncodeToString(id)),
	//} )
	result, err := d.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("storage"),
		Key: aws.String(hex.EncodeToString(id)),
	} )
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
