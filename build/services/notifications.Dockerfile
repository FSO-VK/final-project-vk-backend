FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build ./cmd/notifications/main.go

FROM alpine:3.22

WORKDIR /notifications

COPY --from=builder /build/main /build/config/ ./

ENTRYPOINT ["./main"]

CMD ["--file", "./notifications-conf.yaml"]

EXPOSE 8000