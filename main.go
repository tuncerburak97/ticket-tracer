package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
	"os"
	"sync"
	"ticker-tracer/api"
	"ticker-tracer/config"
	"ticker-tracer/config/db"
	"ticker-tracer/scheduler/train"
	"time"
)

func runScheduler(wg *sync.WaitGroup) {
	defer wg.Done()
	scheduler := train.GetTrainSchedulerInstance()
	err := gocron.Every(15).Seconds().Do(scheduler.Run)
	if err != nil {
		fmt.Printf("Scheduler error: %v\n", err)
		return
	}
	gocron.Start()
}

func runServer(wg *sync.WaitGroup) {
	defer wg.Done()
	config.InitConfig()
	if err := db.InitDb(); err != nil {
		fmt.Printf("Database initialization error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Server is starting...")
	if err := api.InitServer(); err != nil {
		fmt.Printf("Server initialization error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Fatalf("Timezone yüklenemedi: %v", err)
	}
	time.Local = loc
	var wg sync.WaitGroup
	wg.Add(2)

	go runScheduler(&wg)
	go runServer(&wg)

	wg.Wait()
}
