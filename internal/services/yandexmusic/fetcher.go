package yandexmusic

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type FetcherImpl interface {
	Fetch() []int32
}

type Fetch struct {
	newReleases []int32
}

type Releases struct {
	Result struct {
		NewReleases []int32 `json:"newReleases"`
	} `json:"result"`
}

func (n *Fetch) Fetch(token string) Fetch {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.music.yandex.net/landing3/new-releases", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "OAuth "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var releases Releases
	err = json.Unmarshal(body, &releases)
	if err != nil {
		fmt.Println("Can't unmarshal JSON:", err)
	}

	return Fetch{newReleases: releases.Result.NewReleases}
}
