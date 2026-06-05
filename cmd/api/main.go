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

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/docs"
	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	"github.com/rahmatrdn/go-skeleton/internal/http/handler"
	"github.com/rahmatrdn/go-skeleton/internal/parser"
	"github.com/rahmatrdn/go-skeleton/internal/presenter/json"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	comment_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/comment"
	post_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/post"
	todo_list_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list"
	user_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/user"

	"github.com/gofiber/contrib/v3/monitor"
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/subosito/gotenv"
)

func init() {
	_ = gotenv.Load()
}

// @description 				This is a sample swagger for Go Skeleton
// @termsOfService 				http://swagger.io/terms/
// @contact.name 				API Support
// @contact.email 				rahmatrdn.dev@gmail.com
// @license.name				Apache 2.0
// @securityDefinitions.apikey 	Bearer
// @in							header
// @name						Authorization
// @description				    Type "Bearer" followed by a space and JWT token. Example: "Bearer {token}"
// @license.url 				http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	// Initialize config variable from .env file
	cfg := config.NewConfig()
	config.SetTimezone(cfg.AppTimezone)

	// Override Swagger info from env
	docs.SwaggerInfo.Title = cfg.AppName
	docs.SwaggerInfo.Version = cfg.AppVersion
	docs.SwaggerInfo.Host = cfg.ApiHost

	app := fiber.New(config.NewFiberConfiguration(cfg))
	app.Get("/apidoc/*", swaggo.New(swaggo.Config{
		Title:           cfg.AppName + " API Documentation",
		Layout:          "StandaloneLayout",
		URL:             "doc.json",
		DeepLinking:     false,
		DocExpansion:    "none",
		SyntaxHighlight: &swaggo.SyntaxHighlightConfig{Activate: true, Theme: "agate"},
	}))

	app.Get("/metrics", monitor.New(monitor.Config{Title: cfg.AppName + " Metrics Page"}))

	// Middleware setup
	setupMiddleware(app, cfg)

	// logger, _ := config.NewZapLog(cfg.AppEnv)
	// logger = logger.WithOptions(zap.AddCallerSkip(1))

	presenterJson := json.NewJsonPresenter()
	parser := parser.NewParser()

	// RabbitMQ Configuration (if needed)
	// queue, err := config.NewRabbitMQInstance(context.Background(), &cfg.RabbitMQOption)
	// if err != nil {zp
	// 	log.Fatal(err)
	// }

	// Redis Configuration (if needed)
	// redisDB := config.NewRedis(&cfg.RedisOption)

	// MySQL/MariaDB Initialization
	gormLogger := config.NewGormLogMysqlConfig(&cfg.MysqlOption)
	mysqlDB, err := config.NewMysql(cfg.AppEnv, cfg.AppTimezone, &cfg.MysqlOption, gormLogger)
	if err != nil {
		log.Fatal(err)
	}

	// PostgreSQL Initialization
	// gormLogger := config.NewGormLogPostgreConfig(&cfg.MysqlOption)
	// postgreDB, err := config.NewPostgreSQL(cfg.AppEnv, &cfg.PostgreSqlOption, gormLogger)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// AUTH : Write authetincation mechanism method (JWT, Basic Auth, etc.)
	jwtAuth := auth.NewJWTAuth()

	// REPOSITORY : Write repository code here (database, cache, etc.)
	userRepo := mysql.NewUserRepository(mysqlDB)
	todoListRepo := mysql.NewTodoListRepository(mysqlDB)
	postRepo := mysql.NewPostRepository(mysqlDB)
	commentRepo := mysql.NewCommentRepository(mysqlDB)

	// USECASE : Write bussines logic code here (validation, business logic, etc.)
	// _ = usecase.NewLogUsecase(queue)  // LogUsecase is a sample usecase for sending log to queue (Mongodb, ElasticSearch, etc.)
	userUsecase := user_usecase.NewUserUsecase(userRepo, jwtAuth)
	todoListUsecase := todo_list_usecase.NewTodoListUsecase(todoListRepo)
	postUsecase := post_usecase.NewPostUsecase(postRepo)
	commentUsecase := comment_usecase.NewCommentUsecase(commentRepo, postRepo)

	api := app.Group("/api/v1")

	handler.NewAuthHandler(parser, presenterJson, userUsecase).Register(api)
	handler.NewTodoListHandler(parser, presenterJson, todoListUsecase).Register(api)
	handler.NewPostHandler(parser, presenterJson, postUsecase).Register(api)
	handler.NewCommentHandler(parser, presenterJson, commentUsecase).Register(api)

	app.Get("/health-check", healthCheck)

	// Handle Route not found
	app.Use(routeNotFound)

	runServerWithGracefulShutdown(app, cfg.ApiPort, 30)
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Enable CORS if API shared in public
	// if cfg.AppEnv == "production" {
	// 	app.Use(
	// 		cors.New(cors.Config{
	// 			AllowCredentials: true,
	// 			AllowOrigins:     cfg.AllowedCredentialOrigins,
	// 			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	// 			AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
	// 		}),
	// 	)
	// }

	app.Use(
		logger.New(logger.Config{
			Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			TimeZone:   "Asia/Jakarta",
		}),
		recover.New(recover.Config{
			StackTraceHandler: func(c fiber.Ctx, e interface{}) {
				fmt.Println(c.Request().URI())
				stacks := fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack())
				log.Println(stacks)
			},
			EnableStackTrace: true,
		}),
	)
}

func runServerWithGracefulShutdown(app *fiber.App, apiPort string, shutdownTimeout int) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Run server in a goroutine
	go func() {
		defer wg.Done()
		log.Printf("Starting REST server, listening at %s\n", apiPort)
		if err := app.Listen(apiPort); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Capture OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down REST server...")

	// Timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	} else {
		log.Println("REST server shut down gracefully")
	}

	// Wait for goroutines to exit
	wg.Wait()
	log.Println("All tasks completed. Exiting application.")
}

var healthCheck = func(c fiber.Ctx) error {
	return c.JSON(entity.GeneralResponse{
		Code:    200,
		Message: "OK!",
	})
}

var routeNotFound = func(c fiber.Ctx) error {
	return c.Status(404).JSON(entity.GeneralResponse{
		Code:    404,
		Message: "Route Not Found!",
	})
}
