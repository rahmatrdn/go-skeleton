package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/subosito/gotenv"
)

func init() {
	_ = gotenv.Load()
}

func main() {
	location, _ := time.LoadLocation("Asia/Jakarta")

	s, err := gocron.NewScheduler(
		gocron.WithLocation(location),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting scheduler...")

	// cfg := config.NewConfig()
	// queue, err := config.NewRabbitMQInstance(context.Background(), &cfg.RabbitMQOption)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			4*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
				fmt.Println("uwu")

				helper.LogInfo("Process", "func_name", entity.CaptureFields{}, "message")
			},
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
	}

	s.Start()
	fmt.Println("Scheduler started!")

	// Keep the main program running indefinitely
	select {} // Infinite loop
}
