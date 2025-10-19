FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build ./cmd/medication/main.go

FROM alpine:3.22

WORKDIR /medication

COPY --from=builder /build/main /build/config/ ./

ENTRYPOINT ["./main"]

CMD ["--file", "./auth-conf.yaml"]

EXPOSE 8080