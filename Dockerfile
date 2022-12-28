FROM golang:1.19-alpine AS builder

COPY . .
RUN go get -d -v
RUN go build -o /go/bin/rtrouter

FROM scratch
COPY --from=builder /go/bin/rtrouter /go/bin/rtrouter

EXPOSE 8080

ENTRYPOINT ["/go/bin/rtrouter"]