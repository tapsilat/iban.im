FROM golang:1.25-alpine AS builder

RUN apk add --no-cache upx
RUN apk --no-cache add tzdata

WORKDIR /src/ibanim

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ibanim main.go
RUN upx ibanim


FROM scratch

# take env from build args
ARG VERSION
ENV APP_VERSION=$VERSION

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /bin/ibanim

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /src/ibanim/ibanim .


CMD [ "./ibanim" ]