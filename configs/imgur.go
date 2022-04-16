package configs

import "os"

type imgur struct {
	API_BASE_URL string
	CLIENT_ID    string
}

func initImgur() *imgur {
	result := &imgur{
		API_BASE_URL: os.Getenv("IMGUR_API_BASE_URL"),
		CLIENT_ID:    os.Getenv("IMGUR_CLIENT_ID"),
	}

	return result
}
