FROM golang:1.22-alpine

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate swagger docs
RUN swag init -g ./cmd/api/main.go -o ./docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# Run the application
CMD ["./app"] 