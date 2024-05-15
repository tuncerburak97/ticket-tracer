package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"sync"
	"ticker-tracer/api"
	"ticker-tracer/scheduler/train"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := gocron.Every(15).Seconds().Do(
			train.GetTrainSchedulerInstance().Run,
		)
		if err != nil {
			return
		}
		<-gocron.Start()
	}()

	go func() {
		defer wg.Done()
		api.InitServer()
		fmt.Println("Server is running on port 8080")
	}()

	wg.Wait()
}
