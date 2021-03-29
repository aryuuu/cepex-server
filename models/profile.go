package models

import "mime/multipart"

type ProfileUsecase interface {
	UploadPicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type ProfileRepository interface {
}
