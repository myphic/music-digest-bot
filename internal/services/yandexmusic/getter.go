package yandexmusic

import (
	"context"
	"encoding/json"
	"fmt"
	"music-digest-bot/internal/services"
	"net/http"
	"strconv"
	"sync"
)

type resultWithError struct {
	Albums services.Albums
	Err    error
}

func getAlbums(releasesIds []int, wgCount int, token string) ([]services.Albums, error) {
	inputCh := make(chan int)

	outputCh := make(chan resultWithError)
	wg := &sync.WaitGroup{}

	output := make([]services.Albums, 0, len(releasesIds))

	go func() {
		defer close(inputCh)

		for i := range releasesIds {
			inputCh <- releasesIds[i]
		}
	}()

	go func() {
		for i := 0; i < wgCount; i++ {
			wg.Add(1)

			go fetchAlbums(wg, inputCh, outputCh, token)
		}
		wg.Wait()
		close(outputCh)
	}()

	for res := range outputCh {
		if res.Err != nil {
			return nil, fmt.Errorf("an error occurred: %w", res.Err)
		}

		output = append(output, res.Albums)
	}

	return output, nil
}

func fetchAlbums(wg *sync.WaitGroup, inCh <-chan int, outCh chan<- resultWithError, token string) {
	defer wg.Done()
	client := NewClient(&http.Client{})
	for id := range inCh {
		var albums services.Albums

		get, err := client.Get("albums/"+strconv.Itoa(id), token)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(get, &albums)
		if err != nil {
			fmt.Println("Can't unmarshal JSON:", err)
		}
		fmt.Println("albums:", albums)
		outCh <- resultWithError{
			Albums: albums,
			Err:    err,
		}
	}
}

type Yandex struct {
}

func (y Yandex) FetchFromService(ctx context.Context, token string) []services.Albums {
	client := NewClient(&http.Client{})
	body, err := client.Get("landing3/new-releases", token)
	var releases services.Releases
	err = json.Unmarshal(body, &releases)
	if err != nil {
		fmt.Println("Can't unmarshal JSON:", err)
	}
	releasesIds := releases.Result.NewReleases
	albums, err := getAlbums(releasesIds, 5, token)
	fmt.Println(albums)
	if err != nil {
		fmt.Errorf("an error ocurred: %s", err)
	}
	return albums
}
