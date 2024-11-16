package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/internal/queue"
	"github.com/rahmatrdn/go-skeleton/internal/queue/consumer"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mongodb"
	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	_ = gotenv.Load()
}

type GoSkeletonWorker struct {
	ctx     context.Context
	mongoDB *mongo.Database
	queue   queue.Queue
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("[Worker] topic not found, please use 'go run cmd/worker/main.go your.topic-key'")
	}

	log.Println("Starting WORKER")

	var app GoSkeletonWorker
	var err error

	app.ctx = context.Background()
	cfg := config.NewConfig()

	app.mongoDB, err = config.NewMongodb(app.ctx, &cfg.MongodbOption)
	if err != nil {
		log.Fatal(err)
	}
	defer app.mongoDB.Client().Disconnect(app.ctx)

	// gormLogger := config.NewGormLogConfig(&cfg.MysqlOption)
	// mysqlDB, err := config.NewMysql(cfg.AppEnv, &cfg.MysqlOption, gormLogger)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	app.queue, err = config.NewRabbitMQInstance(app.ctx, &cfg.RabbitMQOption)
	if err != nil {
		log.Fatal(err)
	}

	// MongoDB Repository
	logMongoRepo := mongodb.NewLogRepository(app.mongoDB)

	// Consumer
	logConsumer := consumer.NewLogConsumer(context.Background(), logMongoRepo)
	exampleConsumer := consumer.NewExampleConsumer(context.Background(), logMongoRepo)

	var interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	switch os.Args[1] {
	case queue.ProcessSyncLog:
		log.Printf("[Worker] Listening to %v", queue.ProcessSyncLog)
		go app.queue.HandleConsumedDeliveries(queue.ProcessSyncLog, logConsumer.ProcessSyncLog)
	case queue.ProcessExample:
		log.Printf("[Worker] Listening to %v", queue.ProcessExample)
		go app.queue.HandleConsumedDeliveries(queue.ProcessExample, exampleConsumer.Process)
	default:
		log.Fatalf("[Worker] topic not found : %v", os.Args[1])
	}

	<-interrupt
	log.Println("Shutting down the Worker...")

	if err = app.queue.Close(); err != nil {
		log.Printf("Fail shutting down Worker: %s\n", err.Error())
	} else {
		log.Println("Worker successfully shutdown")
	}

}
