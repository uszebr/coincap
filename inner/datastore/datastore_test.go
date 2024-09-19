package datastore

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock response data structure
type MockData struct {
	Value string `json:"value"`
}

// Helper function to create a mock HTTP server
func createMockServer(responseData MockData, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(responseData)
	})
	return httptest.NewServer(handler)
}

func TestNewDataStore(t *testing.T) {
	ds := New[MockData]("mockDataStore", "http://example.com", 1*time.Second, "test-api-key")
	assert.NotNil(t, ds)
	assert.Equal(t, "mockDataStore", ds.name)
	assert.Equal(t, "http://example.com", ds.url)
	assert.Equal(t, 1*time.Second, ds.interval)
	assert.Equal(t, "test-api-key", ds.apiKey)
}

func TestStartPollingAndFetching(t *testing.T) {
	mockData := MockData{Value: "test-value"}
	mockServer := createMockServer(mockData, http.StatusOK)
	defer mockServer.Close()

	// Create a datastore with the mock server URL
	ds := New[MockData]("mockDataStore", mockServer.URL, 500*time.Millisecond, "")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	// Start polling in the background
	go ds.StartPolling(ctx, &wg)

	// Subscribe to the data store
	subChan := ds.Subscribe()

	// Test that data is fetched and published
	select {
	case data := <-subChan:
		assert.Equal(t, "test-value", data.Value)
	case <-time.After(2 * time.Second):
		t.Fatal("Failed to receive data in time")
	}

	// Stop polling
	cancel()
	wg.Wait()
}

func TestFetchAPI(t *testing.T) {
	mockData := MockData{Value: "test-value"}
	mockServer := createMockServer(mockData, http.StatusOK)
	defer mockServer.Close()

	ds := New[MockData]("mockDataStore", mockServer.URL, 1*time.Second, "")

	// Call fetchAPI to get mock data
	data, err := ds.fetchAPI()

	assert.NoError(t, err)
	assert.Equal(t, "test-value", data.Value)
}

func TestFetchAPI_Error(t *testing.T) {
	// Mock server returns a bad request
	mockServer := createMockServer(MockData{}, http.StatusBadRequest)
	defer mockServer.Close()

	ds := New[MockData]("mockDataStore", mockServer.URL, 1*time.Second, "")

	// Call fetchAPI to simulate error
	_, err := ds.fetchAPI()

	assert.Error(t, err, "Expected error when fetching from bad request")
}

func TestUpdateData(t *testing.T) {
	ds := New[MockData]("mockDataStore", "http://example.com", 1*time.Second, "")

	// Update data and verify
	mockData := MockData{Value: "updated-value"}
	ds.updateData(mockData)

	// Ensure data is updated
	assert.Equal(t, "updated-value", ds.GetData().Value)
}

func TestSubscribeAndUnsubscribe(t *testing.T) {
	ds := New[MockData]("mockDataStore", "http://example.com", 1*time.Second, "")

	// Subscribe to the DataStore
	var subChan chan MockData

	subChan = ds.Subscribe()

	// Publish data manually
	mockData := MockData{Value: "test-publish"}
	go func() { ds.publish(mockData) }()

	// Verify the subscriber receives the data
	select {
	case data := <-subChan:
		assert.Equal(t, "test-publish", data.Value)
	case <-time.After(3 * time.Second):
		t.Fatal("Did not receive data in time")
	}

	// Unsubscribe and verify the channel is closed
	ds.Unsubscribe(subChan)
	_, ok := <-subChan
	assert.False(t, ok, "Expected channel to be closed after unsubscribe")
}
