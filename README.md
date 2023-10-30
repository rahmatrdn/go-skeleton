
# Propper Skeleton for your Golang Project

## Description
`go-skeleton` is a boilerplate for Golang projects. The project structure follows the Clean Code Architecture ([Read here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)). This Skeleton made with **Fiber Framework**.  

"Forget about the complexities of folder structures, and concentrate on growing your project!"

Features : 
1. Rest API
2. Worker with RabbitMQ
3. Api Docs with Sweager
4. Fiber Framework
5. Implementation of Unit Test (with Testify and Mockery)
6. Authentication with JWT RS512
7. GRPC Server! (IN PROGRESS)
8. GRPC Server with handle authentication (Soon)
9. Caching with Redis (Soon)
10. Dependency Injection with Google Wire (Soon)

Feel free to contribute to this repository if you'd like!


## Contact
| Name              | Email                           | Role       |
| :----------------:|:-------------------------------:|:----------:|
| Rahmat Ramadhan Putra  | rahmat.putra@spesolution.com    | Creator   |


## Development Guide
### Prerequisite
- Git (See [Git Installation](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git))
- Go 1.19 (See [Golang Installation](https://golang.org/doc/install))
- Go Migrate CLI (See [Migrate CLI Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- MySQL (See [MySQL Installation](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/))
- Redis (See [Redis Installation](https://redis.io/docs/getting-started/installation/))
- Mockery (See [Mockery Installation](https://github.com/vektra/mockery))

#### Windows OS (for a better development experience)

*   Install [Make](https://www.gnu.org/software/make/) (See [Make Installation](https://leangaurav.medium.com/how-to-setup-install-gnu-make-on-windows-324480f1da69)).


### Installation
1. Clone this repo
```sh
git clone <repo_url>
```
2. Copy `example.env` to `.env`
```sh
cp env.sample .env
```
3. Create a MySQL database with the name "go\_skeleton."
4. Run database migration
```sh
make migrate_up
```
- Generate `private_key.pem` and `public_key.pem`. You can generate them using an [Online RSA Generator](https://travistidwell.com/jsencrypt/demo/) or other tools. Place the files in the project's root folder.
- Start the API Service
```sh
go run cmd/api/main.go
```
- Start the Worker Service (if needed)
```sh
go run cmd/worker/main.go
```

### Running In Docker
- Docker Build for API
```sh
docker build -t go-skeleton-api:1.0.1-dev -f ./deploy/docker/api/Dockerfile .
```
- Docker Build for Worker
```sh
docker build -t go-skeleton-worker:1.0.1-dev -f ./deploy/docker/worker/Dockerfile .
```
- Run docker compose for API and Workers
```sh
docker-compose -f docker-compose.yaml -f docker-compose-worker.yaml -f docker-compose-worker-2.yaml up -d
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
make mockery dependency=DependencyClassName
```
- Run unit test with command below or You can run test per function using Vscode!
```sh
make test
```

## Contributing
- Create a new branch with a descriptive name that reflects the changes and switch to the new branch. Use the prefix feature/ for new features or fix/ for bug fixes.
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