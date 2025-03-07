# Builder Stage
FROM golang:1.23 as builder

WORKDIR /usr/app

# Proje dosyalarını kopyala
COPY . .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ibanim .

# Migrations'ı derle
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./migrations/

# Multi-stage build

# Alpine Stage (minimal imaj)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

ENV PATH /root

WORKDIR /root/

# Builder aşamasından derlenmiş ikili dosyaları al
COPY --from=builder /usr/app/ibanim .
COPY --from=builder /usr/app/migrate .

# Konfigürasyon ve şablon dosyalarını kopyala
COPY --from=builder /usr/app/config /root/config
COPY --from=builder /usr/app/templates /root/templates

# Örnek konfigürasyon dosyalarını kopyala
COPY --from=builder /usr/app/config/application.yml /root/config/application.yml
COPY --from=builder /usr/app/config/smtp.yml /root/config/smtp.yml
COPY --from=builder /usr/app/config/database.yml /root/config/database.yml

EXPOSE 8080

# Başlatma komutunu belirt
CMD ["./ibanim"]
