package services

import (
	"context"
	"music-digest-bot/internal/db/repository"
	"sync"
)

type FetcherImpl interface {
	FetchFromService(ctx context.Context, token string) []Albums
}

type Fetch struct {
	sources     repository.SourcesRepositoryImpl
	fetchers    FetcherImpl
	newReleases []int
}

type Releases struct {
	Result struct {
		NewReleases []int `json:"newReleases"`
	} `json:"result"`
}

type Albums struct {
	Result struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Type        string `json:"type"`
		ReleaseDate string `json:"releaseDate"`
		Genre       string `json:"genre"`
		LikesCount  int    `json:"likesCount"`
		Artists     []struct {
			Name string `json:"name"`
		}
	} `json:"result"`
}

func (f *Fetch) Fetch(ctx context.Context, token string) error {
	sources, err := f.sources.Sources(ctx)

	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, source := range sources {
		wg.Add(1)

		go func(source repository.SourceModel) {
			f.fetchers.FetchFromService(ctx, token)
			defer wg.Done()
		}(source)

	}
	wg.Wait()

	return nil
}
