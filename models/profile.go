package models

type ProfileUsecase interface {
	// UploadPicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UploadAvatar(file string) (string, error)
}

type ProfileRepository interface {
}
