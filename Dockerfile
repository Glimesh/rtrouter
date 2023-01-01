FROM golang:1.19-alpine AS build

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /rtrouter

FROM scratch AS runtime

WORKDIR /
EXPOSE 8080

COPY --from=build /rtrouter /rtrouter

ENTRYPOINT ["/rtrouter"]