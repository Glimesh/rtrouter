FROM golang:1.19-alpine AS build

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/rtrouter ./...

ENTRYPOINT ["/usr/local/bin/rtrouter"]