package main

import (
	"coincap/domain/coincapdomain"
	"coincap/inner/datastore"
	"coincap/inner/logger/loggerinit"
	"context"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	//TODO save interval in config for different stores
	//TODO add test, dev, prod in config; Setup logger with config settings
	loggerinit.MustInitLogger(loggerinit.LogDebug)
	if err := godotenv.Load(); err != nil {
		log.Fatal("godotenv issue")
	}
	coincapKey := os.Getenv("COINCAP_KEY")
	if coincapKey == "" {
		slog.Info("COINCAP_KEY env var is not set. Lower quantity of request per minute available")
	}

	hostName := os.Getenv("HOST_NAME")
	if hostName == "" {
		log.Fatal("HOST_NAME env var is not set")
	}

	slog.Debug("coincap started..")

	var wg sync.WaitGroup

	// Create API-specific DataStores and start polling
	api1Store := datastore.New[coincapdomain.CoincapAssetsResponse]("assets", "https://"+hostName+"/v2/assets", 2*time.Second, coincapKey)
	api2Store := datastore.New[coincapdomain.CoincapMarketsResponse]("markets", "https://"+hostName+"/v2/assets/bitcoin/markets", 5*time.Second, coincapKey)
	api3Store := datastore.New[coincapdomain.CoincapAssetIdResponse]("bitcoinprice", "https://"+hostName+"/v2/assets/bitcoin", 1*time.Second, coincapKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Start polling and retrieve PubSub objects
	wg.Add(3)
	api1Store.StartPolling(ctx, &wg)
	api2Store.StartPolling(ctx, &wg)
	api3Store.StartPolling(ctx, &wg)

	//api2PubSub := api2Store.StartPolling(ctx, &wg)

	// Subscribers for each API
	go func() {
		ch := api1Store.Subscribe()
		defer api1Store.Unsubscribe(ch)
		for data := range ch {
			slog.Info("--1--API1: Received data: %+v\n", "size", len(data.Data), "timestamp", data.Timestamp)
		}
	}()
	go func() {
		ch := api2Store.Subscribe()
		defer api2Store.Unsubscribe(ch)
		for data := range ch {
			slog.Info("--2--API2: Received data: %+v\n", "size", len(data.Data), "timestamp", data.Timestamp)
		}
	}()

	go func() {
		ch := api3Store.Subscribe()
		defer api3Store.Unsubscribe(ch)
		for data := range ch {
			slog.Info("--3--API2: Received data: %+v\n", "name", data.Data.Name, "price", data.Data.PriceUsd)
		}
	}()

	// go func() {
	// 	ch := api2PubSub.Subscribe()
	// 	defer api2PubSub.Unsubscribe(ch)
	// 	for data := range ch {
	// 		log.Printf("API2: Received data: %+v\n", data)
	// 	}
	// }()

	// Run for 60 seconds
	time.Sleep(50 * time.Second)
	cancel() // Stop all polling
	wg.Wait()

}
