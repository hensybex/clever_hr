# ./api/Dockerfile

# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the built server binary from the builder stage
COPY --from=builder /app/server .
# Copy the wait-for.sh script into the image
COPY wait-for.sh .

# Install postgresql-client and poppler-utils (pdftotext)
RUN apk add --no-cache postgresql-client poppler-utils

# Make wait-for.sh executable
RUN chmod +x wait-for.sh

EXPOSE 8080

# Command to run the application
CMD ["./wait-for.sh", "db", "./server"]
