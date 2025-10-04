FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk update && apk add --no-cache curl tar 

ADD go.mod .

COPY . .

#RUN go build -o app ./cmd/app/main.go

FROM alpine:latest

#WORKDIR /app

# COPY --from=builder /build/app .

# ENTRYPOINT ["./app"]

EXPOSE 8000
