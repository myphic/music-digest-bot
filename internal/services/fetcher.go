package services

import (
	"context"
	"music-digest-bot/internal/db/repository"
	"sync"
)

type Fetcher interface {
	FetchFromService(ctx context.Context) []Albums
}

type FetchImpl struct {
	sources     repository.SourcesRepository
	digest      repository.DigestRepository
	fetchers    Fetcher
	newReleases []int
}

func New(sources repository.SourcesRepository, digest repository.DigestRepository, fetchers Fetcher) *FetchImpl {
	return &FetchImpl{
		sources:  sources,
		digest:   digest,
		fetchers: fetchers,
	}
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

func (f *FetchImpl) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)

	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, source := range sources {
		wg.Add(1)
		var albums []Albums /* todo save albums to storage */
		go func(source repository.SourceModel, albums []Albums) {
			albums = f.fetchers.FetchFromService(ctx)
			f.processItems(ctx, source, albums)
			defer wg.Done()
		}(source, albums)
	}
	wg.Wait()

	return nil
}

func (f *FetchImpl) processItems(ctx context.Context, source repository.SourceModel, albums []Albums) error {
	return nil
}
