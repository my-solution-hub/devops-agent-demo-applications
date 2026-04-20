# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /app

# Install git for module downloads
RUN apk --no-cache add git

# Use China-friendly Go proxy (remove this line if not needed)
ENV GOPROXY=https://goproxy.cn,direct

# Copy everything and build
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o device-service ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/device-service .
EXPOSE 8082
ENTRYPOINT ["./device-service"]
