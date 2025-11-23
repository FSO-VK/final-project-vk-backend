FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build ./cmd/medication/main.go

FROM alpine:3.22

WORKDIR /medication

COPY --from=builder /build/main /build/config/ ./
COPY --from=builder /build/templates/ ./templates/
COPY --from=builder /build/internal/medication/infrastructure/llm_chat_bot/templates/ ./internal/medication/infrastructure/llm_chat_bot/templates/

ENTRYPOINT ["./main"]

CMD ["--file", "./medication-conf.yaml"]

EXPOSE 8080