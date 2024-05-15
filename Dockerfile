# Build aşaması
FROM golang:1.21 AS builder

# Çalışma dizinini ayarla
WORKDIR /app

# go mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./

# Bağımlılıkları indir
RUN go mod download

# Kaynak kodu konteynere kopyala
COPY . .

# Uygulamayı derle
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Çalıştırma aşaması
FROM debian:buster-slim

# Çalışma dizinini ayarla
WORKDIR /app

# Derlenmiş binary'yi önceki aşamadan kopyala
COPY --from=builder /app/main .

# 8080 portunu dışa aç
EXPOSE 8080

# Çalıştır komutunu belirt
CMD ["./main"]