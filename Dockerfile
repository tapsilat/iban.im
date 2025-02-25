FROM golang:1.20 AS builder

WORKDIR /usr/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ibanim .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./migrations/

# Multi stage

FROM alpine:latest

RUN apk --no-cache add ca-certificates

ENV PATH=/root:$PATH

WORKDIR /root/

COPY --from=builder /usr/app/ibanim .

COPY --from=builder /usr/app/migrate .

EXPOSE 8080

CMD ["./ibanim"] 
