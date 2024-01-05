package main

import (
	"context"
	"log"
	"time"
)

type Item struct {
	Name string
}

type Subscription interface {
	Updates() <-chan Item
}

type Fetcher interface {
	Fetch() (Item, error)
}

type fetcher struct {
	url string
}

type subscription struct {
	fetcher Fetcher
	updates chan Item
}

func (f *fetcher) Fetch() (Item, error) {
	return Item{Name: "test"}, nil
}

func (s *subscription) Updates() <-chan Item {
	return s.updates
}

func (s *subscription) serve(ctx context.Context, freq int) {
	clock := time.NewTicker(time.Duration(freq) * time.Second)

	type fetchResult struct {
		item Item
		err  error
	}

	fetchDone := make(chan fetchResult, 1)

	for {
		select {
		case <-clock.C:
			// trigger fetch
			go func() {
				item, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{item, err}
			}()
		case result := <-fetchDone:
			if result.err != nil {
				log.Println("fetch error: ", result.err)
				break
			}

			s.updates <- result.item
		case <-ctx.Done():
			return
		}
	}
}

func NewSubscription(ctx context.Context, fetcher Fetcher, freq int) Subscription {
	s := &subscription{
		fetcher: fetcher,
		updates: make(chan Item),
	}

	go s.serve(ctx, freq)
	return s
}

func NewFetcher(url string) Fetcher {
	f := &fetcher{
		url: url,
	}

	return f
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	subscription := NewSubscription(ctx, NewFetcher("http://example.com"), 3)

	time.AfterFunc(10*time.Second, func() {
		cancel()
		log.Println("canceled")
	})

	for item := range subscription.Updates() {
		log.Println(item)
	}
}
