package datastore

import (
	"coincap/inner/logger/logutil"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type DataStore[T any] struct {
	name        string
	url         string
	interval    time.Duration
	datamutex   *sync.RWMutex
	data        T
	subscribers map[chan T]struct{}
	subMutex    *sync.RWMutex
	apiKey      string
}

// Creating DataStore
func New[T any](name string, url string, interval time.Duration, apiKey string) *DataStore[T] {
	return &DataStore[T]{
		name:        name,
		url:         url,
		interval:    interval,
		datamutex:   &sync.RWMutex{},
		data:        *new(T),
		subscribers: make(map[chan T]struct{}),
		subMutex:    &sync.RWMutex{},
		apiKey:      apiKey,
	}
}

// StartPolling begins polling the API at regular intervals, updating data for non subscribers; publish data for subscribers
func (ds *DataStore[T]) StartPolling(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	slog.Debug("starting polling data..")
	go func() {
		ticker := time.NewTicker(ds.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				slog.Debug("Stopp polling ", "url", ds.url)
				return
			case <-ticker.C:
				data, err := ds.fetchAPI()
				if err != nil {
					slog.Error("Error fetching data ", "url", ds.url, logutil.Err(err))
					continue
				}
				ds.updateData(data)
				ds.publish(data) // Automatically publish new data to subscribers like own websocket connection
			}
		}
	}()
}

func (ds *DataStore[T]) updateData(d T) {
	ds.datamutex.Lock()
	defer ds.datamutex.Unlock()
	ds.data = d
}

// Getter for entities that not useing Publish/Subscribe model like own rest api
func (ds *DataStore[T]) GetData() T {
	ds.datamutex.RLock()
	defer ds.datamutex.RUnlock()
	return ds.data
}

// fetchAPI fetches the raw data from the API and unmarshals it into the type T
func (ds *DataStore[T]) fetchAPI() (T, error) {
	slog.Debug("fetching data from api")
	client := &http.Client{}
	req, err := http.NewRequest("GET", ds.url, nil)
	if err != nil {
		return *new(T), err
	}

	if ds.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+ds.apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return *new(T), err
	}
	defer resp.Body.Close()
	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return *new(T), fmt.Errorf("received non-2xx response: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return *new(T), err
	}
	return result, nil
}

func (ds *DataStore[T]) Subscribe() chan T {
	ch := make(chan T)
	ds.subMutex.Lock()
	ds.subscribers[ch] = struct{}{}
	ds.subMutex.Unlock()
	slog.Debug("Subscribed a new channel")
	return ch
}

func (ds *DataStore[T]) Unsubscribe(ch chan T) {
	ds.subMutex.Lock()
	defer ds.subMutex.Unlock()
	if _, ok := ds.subscribers[ch]; ok {
		delete(ds.subscribers, ch)
		close(ch) // Ensure the channel is closed
	}
	slog.Debug("Unsubscribed a channel")
}

func (ds *DataStore[T]) publish(data T) {
	ds.subMutex.RLock()
	defer ds.subMutex.RUnlock()
	for sub := range ds.subscribers {
		select {
		case sub <- data:
		default:
			// Avoid blocking if the subscriber is not ready to receive
			slog.Error("subscriber not ready to receive the data")
		}
	}
	slog.Debug("Published data to all subscribers")
}
