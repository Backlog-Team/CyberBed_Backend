# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine3.18 AS BUILDER

RUN go version

COPY . /github.com/cyber_bed/
WORKDIR /github.com/cyber_bed/

RUN go mod download
RUN GOOS=linux go build -o ./bin/server ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=BUILDER /github.com/cyber_bed/bin/server .
COPY --from=BUILDER /github.com/cyber_bed/configs/ configs/

EXPOSE 8080

CMD ["./server", "-ConfigPath", "configs/app/deploy.yaml"]
