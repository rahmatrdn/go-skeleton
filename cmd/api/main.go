package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"github.com/rahmatrdn/go-skeleton/config"
	_ "github.com/rahmatrdn/go-skeleton/docs"
	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	"github.com/rahmatrdn/go-skeleton/internal/http/handler"
	"github.com/rahmatrdn/go-skeleton/internal/parser"
	"github.com/rahmatrdn/go-skeleton/internal/presenter/json"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/subosito/gotenv"
	glogger "gorm.io/gorm/logger"
)

func init() {
	_ = gotenv.Load()
}

// @title 						Go Skeleton!
// @version 					1.0
// @description 				This is a sample swagger for Go Skeleton
// @termsOfService 				http://swagger.io/terms/
// @contact.name 				API Support
// @contact.email 				rahmat.putra@spesolution.com
// @license.name				Apache 2.0
// @securityDefinitions.apikey 	Bearer
// @in							header
// @name						Authorization
// @license.url 				http://www.apache.org/licenses/LICENSE-2.0.html
// @host 						localhost:7011
// @BasePath /
func main() {
	// Initialize config variable from .env file
	cfg := config.NewConfig()

	app := fiber.New(config.NewFiberConfiguration())

	// Initialize Swagger for API documentation
	app.Get("/apidoc/*", swagger.HandlerDefault)

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     cfg.AllowedCredentialOrigins[0],
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
	}), logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}), recover.New(recover.Config{
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			fmt.Println(c.Request().URI())
			stacks := fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack())
			log.Println(stacks)
		},
		EnableStackTrace: true,
	}))

	presenterJson := json.NewJsonPresenter()
	parser := parser.NewParser()

	// RabbitMQ Configuration (if needed)
	// queue, err := config.NewRabbitMQInstance(context.Background(), &cfg.RabbitMQOption)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	// Redis Configuration (if needed)
	// redis := config.NewRedis(config)

	// AUTH : Write authetincation mechanism method (JWT, Basic Auth, etc.)
	jwtAuth := auth.NewJWTAuth()

	// REPOSITORY : Write repository code here (database, cache, etc.)
	userRepo := mysql.NewUserRepository(mysqlDB)
	todoListRepo := mysql.NewTodoListRepository(mysqlDB)

	// USECASE : Write bussines logic code here (validation, business logic, etc.)
	// _ = usecase.NewLogUsecase(queue)  // LogUsecase is a sample usecase for sending log to queue (Mongodb, ElasticSearch, etc.)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtAuth)
	todoListUsecase := usecase.NewTodoListUsecase(todoListRepo)

	api := app.Group("/v1/api")

	handler.NewAuthHandler(parser, presenterJson, userUsecase).Register(api)
	handler.NewTodoListHandler(parser, presenterJson, todoListUsecase).Register(api)

	app.Get("/health-check", healthCheck)
	app.Get("/metrics", monitor.New())

	var wg = sync.WaitGroup{}
	wg.Add(1)

	// Running server in Goroutines
	go func() {
		defer wg.Done()

		log.Printf("Starting REST, listening at %s\n", cfg.ApiPort)

		if err := app.Listen(cfg.ApiPort); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the REST server...")

	_, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Printf("Fail shutting down REST server: %s\n", err.Error())

		log.Fatal(err)
	}

	log.Println("REST server successfully shutdown")
	wg.Wait()
}

var healthCheck = func(c *fiber.Ctx) error {
	return c.JSON(entity.GeneralResponse{
		Code:    200,
		Message: "OK!",
	})
}
