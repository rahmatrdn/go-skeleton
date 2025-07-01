
# Propper Skeleton for your Golang Project

## Description
`go-skeleton` is a boilerplate for Golang projects. The project structure follows the Clean Code Architecture ([Read here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)). This Skeleton made with **Fiber Framework**.  

"Forget about the complexities of folder structures in Go, focus on your project!"

Principles :
- Reusable and Maintainable Code
- Decoupled Code
- Scalable Development

Features : 
1. REST API
2. Clean Architecture
3. Fiber Framework
4. Api Docs with Swagger
5. Worker/Consumer Queue with RabbitMQ
6. Implementation of Unit Test (with Testify and Mockery)
7. Authentication with JWT RS512
8. Logging with Zap Log
9. GRPC Server! (IN PROGRESS)
10. GRPC Server with handle authentication (Soon!)
11. Caching with Redis (Soon!)
12. Dependency Injection with Google Wire (Soon!)
13. Worker Queue with Kafka (Soon!)

Feel free to contribute to this repository if you'd like!


## Contact
| Name                   | Email                        | Role    |
| ---------------------- | ---------------------------- | ------- |
| Rahmat Ramadhan Putra  | rahmatrdn.dev@gmail.com     | Creator |



## Development Guide
### Prerequisite
- Git (See [Git Installation](https://git-scm.com/downloads))
- Go 1.24+ (See [Golang Installation](https://golang.org/doc/install))
- MySQL / MariaDB / PostgreSQL (Download via Docker or Other sources)
- Mockery (Optional) (See [Mockery Installation](https://github.com/vektra/mockery))
- Go Migrate CLI (Optional) (See [Migrate CLI Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- Redis (Optional based on your requirement) (See [Redis Installation](https://redis.io/docs/getting-started/installation/) or use in Docker)
- RabbitMQ (Optional based on your requirement) (See [RabbitMQ Installation](https://www.rabbitmq.com/download.html) or use in Docker)

#### Windows OS (for a better development experience)

*   Install [Make](https://www.gnu.org/software/make/) (See [Make Installation](https://leangaurav.medium.com/how-to-setup-install-gnu-make-on-windows-324480f1da69)).


### Installation
1. Clone this repo
```sh
git clone https://github.com/rahmatrdn/go-skeleton.git
```
2. Copy `example.env` to `.env`
```sh
cp .env.example .env
```
3. Adjust the `.env` file according to the configuration in your local environment, such as the database or other settings 
4. Create a MySQL database with the name `go_skeleton`
5. Run database migration or Manually run in you SQL Client
```sh
make migrate_up
```
6. Generate `private_key.pem` and `public_key.pem`. You can generate them using an [Online RSA Generator](https://travistidwell.com/jsencrypt/demo/) or other tools. Place the files in the project's root folder.
7. Start the API Service
```sh
go run cmd/api/main.go
```
8. Start the Worker Service (if needed)
```sh
go run cmd/worker/main.go
```

### Api Documentation
For API docs, we are using [Swagger](https://swagger.io/) with [Swag](https://github.com/swaggo/swag) Generator
- Install Swag
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
- Generate apidoc
```sh
make apidoc
```
- Start API documentations
```sh
go run cmd/api/main.go
```
- Access API Documentation with  browser http://localhost:PORT/apidoc



### Unit test
*tips: if you use `VS Code` as your code editor, you can install extension `golang.go` and follow tutorial [showing code coverage after saving your code](https://dev.to/vuong/golang-in-vscode-show-code-coverage-of-after-saving-test-8g0) to help you create unit test*

- Use [Mockery](https://github.com/vektra/mockery) to generate mock class(es)
```sh
make mock d=DependencyClassName
```
- Run unit test with command below or You can run test per function using Vscode!
```sh
make test
```


### Running In Docker
- Docker Build for API
```sh
docker build -t go-skeleton-api:1.0.1 -f ./deploy/docker/api/Dockerfile .
```
- Docker Build for Worker
```sh
docker build -t go-skeleton-worker:1.0.1 -f ./deploy/docker/worker/Dockerfile .
```
- Run docker compose for API and Workers
```sh
docker-compose -f docker-compose.yaml up -d
```


## Contributing
- Create a new branch with a descriptive name that reflects the changes and switch to the new branch. Use the prefix `feature/` for new features or `fix/` for bug fixes.
```sh
git checkout -b <prefix>/branch-name
```
- Make your change(s) and make the test(s)
- Commit and push your change to upstream repository
```sh
git commit -m "[Type] a meaningful commit message"
git push origin branch-name
```
- Open Merge Request in Repository (Reviewer Check Contact Info)
- Merge Request will be merged only if review phase is passed.

## More Details Information
Contact Creator!
