
# Skeleton for Golang Project in SPE


## Description

`go-skeleton` is a boilerplate for Golang projects. The project structure follows the Clean Code Architecture ([Read here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)). Feel free to contribute to this repository if you'd like!

## Owner

Tribe - Ursa Major

## Contact
| Name              | Email                           | Role       |
| :----------------:|:-------------------------------:|:----------:|
| Rahmat Ramadhan  | rahmat.putra@spesolution.com    | Inisiator   |


## Development Guide
### Prerequisite

- Git (See [Git Installation](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git))
- Go 1.19 (See [Golang Installation](https://golang.org/doc/install))
- Go Migrate CLI (See [Migrate CLI Installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- MySQL (See [MySQL Installation](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/))
- Redis (See [Redis Installation](https://redis.io/docs/getting-started/installation/))
- Mockery (See [Mockery Installation](https://github.com/vektra/mockery))

### Installation
- Clone this repo
```sh
git clone <repo_url>
```
- Copy `example.env` to `.env`
```sh
cp env.sample .env
```
- Running Database Migration
- Generate private_key.pem and public_key.pem. You can generate via [Online RSA Generator](https://travistidwell.com/jsencrypt/demo/) or other tools
- Start API Service
```sh
go run cmd/api/main.go
```
- Start Worker Service (if any)
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

go run cmd/apidoc/main.go

```
- Access API Documentation with [browser](http://localhost:<PORT>/docs/index.html)

### Unit test
*tips: if you use `VS Code` as your code editor, you can install extension `golang.go` and follow tutorial [showing code coverage after saving your code](https://dev.to/vuong/golang-in-vscode-show-code-coverage-of-after-saving-test-8g0) to help you create unit test*

- Use [Mockery](https://github.com/vektra/mockery) to generate mock class(es)
```sh
mockery --name=interfaceName --recursive=true --output=tests/mocks
```
- Run unit test with command below or You can run test per function using Vscode!
```sh
make test
```

## Contributing
- Make new branch with descriptive name about the change(s) and checkout to the new branch with prefix `feature/` or `fix/`
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
Access this document below for detail information about go skeleton
[Techinical Guide - GO Skeleton](https://docs.google.com/document/d/1R0xKkQz-iK7Bovx-F9cZkJFUo6yVtAy1j9L8QUeM2tk/edit?usp=sharing )