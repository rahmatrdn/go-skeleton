package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/rpc/handler"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	wallet "github.com/rahmatrdn/go-skeleton/proto/pb/proto"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc"
	glogger "gorm.io/gorm/logger"
)

func init() {
	_ = gotenv.Load()
}

func main() {
	cfg := config.NewConfig()
	server := setupServer(cfg)

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Starting gRPC server, listening at %s\n", cfg.ApiRpcPort)

		if err := server.Serve(getListener(cfg.ApiRpcPort)); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for termination signals to gracefully stop the server
	waitForShutdown(server)
}

func setupServer(cfg *config.Config) *grpc.Server {
	serverRegister := grpc.NewServer()

	queue, err := config.NewRabbitMQInstance(context.Background(), &cfg.RabbitMQOption)
	if err != nil {
		log.Fatal(err)
	}

	mysqlDBLogger := glogger.New(
		log.New(
			os.Stdout,
			"\r\n",
			log.LstdFlags,
		),
		glogger.Config{
			SlowThreshold:             200 * time.Second,
			LogLevel:                  glogger.Warn,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
		},
	)

	mysqlDB, err := config.NewMysql(cfg.AppEnv, &cfg.MysqlOption, mysqlDBLogger)
	if err != nil {
		log.Fatal(err)
	}

	walletRepo := mysql.NewWalletRepository(mysqlDB)

	logUsecase := usecase.NewLogUsecase(queue)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, logUsecase)

	walletService := handler.NewWalletServiceHandler(walletUsecase)

	// Register Request Handler
	wallet.RegisterWalletServiceServer(serverRegister, walletService)

	return serverRegister
}

func getListener(rpcPort string) net.Listener {
	listener, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	return listener
}

func waitForShutdown(server *grpc.Server) {
	// Create a channel to receive termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until a termination signal is received
	<-stop

	// Gracefully stop the server
	server.GracefulStop()

	log.Println("Server stopped")
}
