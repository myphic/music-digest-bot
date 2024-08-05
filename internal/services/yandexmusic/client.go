package yandexmusic

import (
	"io"
	"log"
	"net/http"
)

type YaHttpClient struct {
	client *http.Client
}

func NewClient(httpClient *http.Client) *YaHttpClient {
	return &YaHttpClient{
		client: httpClient,
	}
}

type ClientImpl interface {
	Get(path string, token string) ([]byte, error)
}

var yandexMusicBaseUrl = "https://api.music.yandex.net/"

func (c *YaHttpClient) Get(path string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", yandexMusicBaseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
