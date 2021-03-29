package configs

import "os"

type s3 struct {
	ACCESS_KEY string
	SECRET_KEY string
	ENDPOINT   string
	BUCKET     string
}

func initS3() *s3 {
	result := &s3{
		ENDPOINT:   os.Getenv("S3_ENDPOINT"),
		BUCKET:     os.Getenv("S3_BUCKET"),
		ACCESS_KEY: os.Getenv("S3_ACCESS_KEY"),
		SECRET_KEY: os.Getenv("S3_SECRET_KEY"),
	}

	return result
}
