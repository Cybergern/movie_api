FROM golang:1.8.5-jessie AS builder
WORKDIR /go
ADD src src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main src/api-test/main.go

FROM alpine:3.7
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root
COPY --from=builder /go/main .
CMD ["./main"]