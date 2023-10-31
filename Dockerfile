FROM golang:1.21.3-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates && apk add --no-cache gcc musl-dev

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -o main ./cmd/api

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main ./main
CMD ["./main"]

