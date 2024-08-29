package oss

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOss struct {
	client *oss.Client
}

func NewAliyunOss(client *oss.Client) *AliyunOss {
	return &AliyunOss{
		client: client,
	}
}

func (oss *AliyunOss) UploadFile(ctx context.Context, fileName string, fileBytes []byte) error {
	bucket, err := oss.client.Bucket("ceit")
	if err != nil {
		return fmt.Errorf("oss get bucket err: %v", err)
	}
	err = bucket.PutObject(fileName, bytes.NewReader(fileBytes))
	if err != nil {
		return fmt.Errorf("oss put file err: %v", err)
	}

	return nil
}
