# Dockerfile References: https://docs.docker.com/engine/reference/builder/
FROM golang:1.23.4-alpine3.21 as builder

WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

#This will copy all the files in our repo to the inside the container at root location.
COPY . . 

# RUN go get -d -v
RUN apk --no-cache add ca-certificates
#build our binary at root location. Binary name will be main. We are using go modules so gpath env variable should be empty.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o deploy/worker ./cmd/worker

#This will use scratch image which is the smallest image we can have. We will use scratch because we needed go environment only for building.
FROM scratch

#we copy our binary from build to scratch.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/deploy/worker .

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

#we tell docker what to run when this image is run and run it as executable.
ENTRYPOINT ["/worker"]