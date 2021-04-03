package usecases

import (
	"mime/multipart"

	"github.com/aryuuu/cepex-server/models"
)

type profileUsecase struct {
	s3Repo models.S3Repository
}

func NewProfileUsecase(s3r models.S3Repository) models.ProfileUsecase {
	return &profileUsecase{
		s3Repo: s3r,
	}
}

func (u *profileUsecase) UploadPicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return u.s3Repo.UploadImage(file, fileHeader)
}
