package oss

import "context"

type Client interface {
	UploadFile(ctx context.Context, fileName string, fileBytes []byte) error
}
