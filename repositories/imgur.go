package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/aryuuu/cepex-server/configs"
	"github.com/aryuuu/cepex-server/models"
)

type imgurRepo struct {
	client *http.Client
}

func NewImgurRepo(client *http.Client) models.ImageRepository {
	result := &imgurRepo{
		client: client,
	}

	return result
}

func (ir *imgurRepo) UploadImageBase64(fileBase64 string) (string, error) {
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("type", "base64")
	_ = writer.WriteField("image", fileBase64)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, configs.Imgur.API_BASE_URL+"/upload", payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer Client-ID "+configs.Imgur.CLIENT_ID)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to upload image: status code: %d", res.StatusCode)
	}
	log.Printf("response status: %d", res.StatusCode)
	defer res.Body.Close()

	var response struct {
		Data models.UploadImageImgurResp `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("imgurRepo.UploadImageBase64: failed to decode response: %v", err)
	}

	return response.Data.Link, nil
}

func (ir *imgurRepo) UploadImageURL(url string) (string, error) {
	return "", nil
}
