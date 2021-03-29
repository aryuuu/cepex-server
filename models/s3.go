package models

import "mime/multipart"

type S3Repository interface {
	UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}
