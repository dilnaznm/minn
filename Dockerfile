
FROM golang:1.23-alpine AS builder


RUN apk add --no-cache git


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o s3uploader main.go


FROM alpine:latest


RUN apk --no-cache add ca-certificates

WORKDIR /root/


COPY --from=builder /app/s3uploader .


ENV S3_ENDPOINT=minio:9000
ENV S3_ACCESS_KEY=minioadmin
ENV S3_SECRET_KEY=minioadmin
ENV S3_USE_SSL=false


CMD ["./s3uploader"]

