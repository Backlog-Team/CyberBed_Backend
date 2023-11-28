# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine3.18 AS BUILDER

RUN go version

COPY . /github.com/cyber_bed/
WORKDIR /github.com/cyber_bed/

RUN go mod download
RUN GOOS=linux go build -o ./bin/notifications ./cmd/notifications/notifications.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/cyber_bed/bin/notifications .
COPY --from=builder /github.com/cyber_bed/configs/ configs/

CMD ["./notifications", "-ConfigPath", "configs/app/notifications/deploy.yaml"]
