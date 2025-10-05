FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build ./cmd/auth/main.go

FROM alpine:3.22

WORKDIR /auth

COPY --from=builder /build/main .

ENTRYPOINT ["./main"]

EXPOSE 8080