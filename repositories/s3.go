package repositories

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/aryuuu/cepex-server/configs"
	"github.com/aryuuu/cepex-server/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/mgo.v2/bson"
)

type s3Repo struct {
	session *session.Session
}

func NewS3Repo(session *session.Session) models.S3Repository {
	result := &s3Repo{
		session: session,
	}

	return result
}

func (m *s3Repo) UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	tempFileName := "pictures/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	_, err := s3.New(m.session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(configs.S3.BUCKET),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType((buffer))),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}

	finalFileName := "https://" + configs.S3.BUCKET + "." + configs.S3.ENDPOINT + "/" + tempFileName
	return finalFileName, err
}
