package usecases

import (
	"github.com/aryuuu/cepex-server/models"
)

type profileUsecase struct {
	// s3Repo    models.S3Repository
	imageRepo models.ImageRepository
}

func NewProfileUsecase(ir models.ImageRepository) models.ProfileUsecase {
	return &profileUsecase{
		// s3Repo:    s3r,
		imageRepo: ir,
	}
}

// func (u *profileUsecase) UploadPicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
// 	return u.s3Repo.UploadImage(file, fileHeader)
// }

func (u *profileUsecase) UploadAvatar(file string) (string, error) {
	return u.imageRepo.UploadImageBase64(file)
}
