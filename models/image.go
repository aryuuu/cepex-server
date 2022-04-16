package models

import "mime/multipart"

type S3Repository interface {
	UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type ImageRepository interface {
	UploadImageURL(url string) (string, error)
	UploadImageBase64(fileBase64 string) (string, error)
}

type UploadImageImgurResp struct {
	Link string `json:"link"`
}
